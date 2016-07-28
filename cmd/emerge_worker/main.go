package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeepers-creepers/emerge/internal/sms"
	"github.com/jeepers-creepers/emerge/internal/subscription"

	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"
)

var (
	pool *pgx.ConnPool
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	cfg, err := pgx.ParseURI(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   cfg,
		AfterConnect: afterConnect,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	client := que.NewClient(pool)

	wm := que.WorkMap{
		subscription.Subscribe: subscribe(pool),
		sms.Inbox:              inbox(pool),
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

func inbox(pool *pgx.ConnPool) que.WorkFunc {
	i := sms.NewInbox(pool)
	return func(job *que.Job) error {
		var m sms.Message
		if err := json.Unmarshal(job.Args, &m); err != nil {
			return err
		}
		if err := i.Add(m); err != nil {
			return err
		}
		return nil
	}
}

func afterConnect(conn *pgx.Conn) (err error) {
	var xs = []func(*pgx.Conn) error{
		que.PrepareStatements,
		subscription.PrepareStatements,
	}
	for _, x := range xs {
		if err = x(conn); err != nil {
			return err
		}
	}
	return
}
