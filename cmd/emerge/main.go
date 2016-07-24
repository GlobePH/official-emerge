package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jeepers-creepers/emerge/internal/channel"
	"github.com/jeepers-creepers/emerge/internal/notify"
	"github.com/jeepers-creepers/emerge/internal/subscribers"
	"github.com/jeepers-creepers/emerge/internal/subscription"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	chain := alice.New(logging, recovery, cors)
	mux := mux.NewRouter().StrictSlash(true)

	apiMux := mux.PathPrefix("/api/").Subrouter()
	ss := subscribers.New(db)
	apiMux.Handle("/subscription", chain.Then(subscription.Handler(ss)))
	apiMux.Handle("/notify", chain.Then(notify.Handler(db)))
	apiMux.Handle("/channel", chain.Then(channel.Handler(db)))

	mux.PathPrefix("/").Handler(chain.Then(http.FileServer(http.Dir("public"))))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
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
