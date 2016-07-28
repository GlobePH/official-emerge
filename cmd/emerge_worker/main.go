package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeepers-creepers/emerge/internal/helpers"
	"github.com/jeepers-creepers/emerge/internal/sms"
	"github.com/jeepers-creepers/emerge/internal/subscription"

	"github.com/bgentry/que-go"
	"github.com/garyburd/redigo/redis"
	"github.com/jackc/pgx"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("$REDIS_URL must be set")
	}

	pgxPool, err := helpers.NewPgxPool(dbURL, afterConnect)
	if err != nil {
		log.Fatal(err)
	}
	defer pgxPool.Close()

	redisPool, err := helpers.NewRedisPool(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	defer redisPool.Close()

	client := que.NewClient(pgxPool)

	wm := que.WorkMap{
		subscription.Subscribe: subscribe(pgxPool),
		sms.Inbox:              inbox(pgxPool, redisPool),
	}

	cSig := make(chan os.Signal)
	signal.Notify(cSig, syscall.SIGTERM, syscall.SIGINT)
	workers := que.NewWorkerPool(client, wm, 2)
	go workers.Start()
	log.Printf("%s started", os.Args[0])
	sig := <-cSig
	log.Printf("%s signal received. Shutting down...", sig)
	workers.Shutdown()
}

func subscribe(pool *pgx.ConnPool) que.WorkFunc {
	ss := subscription.NewSubscribers(pool)
	return func(job *que.Job) error {
		var s subscription.Subscriber
		if err := json.Unmarshal(job.Args, &s); err != nil {
			return err
		}
		if err := ss.Add(s); err != nil {
			return err
		}
		return nil
	}
}

func inbox(pgxPool *pgx.ConnPool, redisPool *redis.Pool) que.WorkFunc {
	return func(job *que.Job) error {
		var m sms.Message
		if err := json.Unmarshal(job.Args, &m); err != nil {
			return err
		}

		i := sms.NewInbox(pgxPool)
		if err := i.Add(m); err != nil {
			return err
		}

		rc := redisPool.Get()
		defer rc.Close()
		if err := rc.Send("PUBLISH", job.Type, job.Args); err != nil {
			return err
		}
		if err := rc.Flush(); err != nil {
			return err
		}
		return nil
	}
}

func afterConnect(conn *pgx.Conn) (err error) {
	var xs = []func(*pgx.Conn) error{
		que.PrepareStatements,
		sms.PrepareStatements,
		subscription.PrepareStatements,
	}
	for _, x := range xs {
		if err = x(conn); err != nil {
			return err
		}
	}
	return
}
