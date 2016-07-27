package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jeepers-creepers/emerge/internal/channel"
	//"github.com/jeepers-creepers/emerge/internal/notify"
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

	chain := alice.New(logging, recovery, cors)
	mux := mux.NewRouter().StrictSlash(true)

	apiMux := mux.PathPrefix("/api/").Subrouter()
	sub := subscription.New(pool)
	apiMux.Handle("/subscription", chain.Then(sub))
	//apiMux.Handle("/notify", chain.Then(notify.Handler(db)))
	apiMux.Handle("/channel", chain.Then(channel.Handler()))

	mux.PathPrefix("/").Handler(chain.Then(http.FileServer(http.Dir("public"))))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("%s started", os.Args[0])
	log.Fatal(s.ListenAndServe())
}

func afterConnect(conn *pgx.Conn) error {
	return que.PrepareStatements(conn)
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
