package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jeepers-creepers/emerge/internal/subscription"

	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"
)

const (
	subscribeSQL = `INSERT INTO subscribers (subscriber_number, access_token) VALUES ($1, $2);`
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
	return func(job *que.Job) error {
		var sub subscription.Subscriber
		if err := json.Unmarshal(job.Args, &sub); err != nil {
			return err
		}
		if _, err := pool.Exec(subscription.Subscribe, sub.SubscriberNumber, sub.AccessToken); err != nil {
			return err
		}
		log.Printf("%s subscribed", sub.SubscriberNumber)
		return nil
	}
}

func afterConnect(conn *pgx.Conn) error {
	if _, err := conn.Prepare(subscription.Subscribe, subscribeSQL); err != nil {
		return err
	}
	return que.PrepareStatements(conn)
}
