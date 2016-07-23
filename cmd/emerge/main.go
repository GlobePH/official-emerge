package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	chain := alice.New(logging, recovery)
	mux := mux.NewRouter().StrictSlash(true)
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
