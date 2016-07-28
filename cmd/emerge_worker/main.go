package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	pgxPool, err := newPgxPool(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgxPool.Close()

	redisPool, err := newRedisPool(redisURL)
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
	i := sms.NewInbox(pgxPool)
	rc := redisPool.Get()
	return func(job *que.Job) error {
		defer rc.Close()
		var m sms.Message
		if err := json.Unmarshal(job.Args, &m); err != nil {
			return err
		}
		if err := i.Add(m); err != nil {
			return err
		}
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

func newRedisPool(uri string) (*redis.Pool, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	var password string
	if u.User != nil {
		password, _ = u.User.Password()
	}
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", u.Host)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}

func newPgxPool(uri string) (pool *pgx.ConnPool, err error) {
	cfg, err := pgx.ParseURI(uri)
	if err != nil {
		return
	}
	return pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   cfg,
		AfterConnect: afterConnect,
	})
}
