package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jeepers-creepers/emerge/internal/channel"
	"github.com/jeepers-creepers/emerge/internal/helpers"
	"github.com/jeepers-creepers/emerge/internal/notify"
	"github.com/jeepers-creepers/emerge/internal/subscription"

	"github.com/bgentry/que-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/justinas/alice"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

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

	chain := alice.New(logging, recovery, cors)
	mux := mux.NewRouter().StrictSlash(true)

	apiMux := mux.PathPrefix("/api/").Subrouter()
	apiMux.Handle("/subscription", chain.Then(subscription.New(pgxPool)))
	apiMux.Handle("/notify", chain.Then(notify.New(pgxPool)))
	apiMux.Handle("/channel", chain.Then(channel.New(redisPool)))

	mux.PathPrefix("/").Handler(chain.Then(http.FileServer(http.Dir("public"))))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("%s started", os.Args[0])
	log.Fatal(s.ListenAndServe())
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

func logging(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}

func recovery(h http.Handler) http.Handler {
	return handlers.RecoveryHandler()(h)
}

func cors(h http.Handler) http.Handler {
	return handlers.CORS()(h)
}
