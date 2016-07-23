package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/jeepers-creepers/emerge/internal/subscription"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	chain := alice.New(debug, logging, recovery)
	mux := mux.NewRouter().StrictSlash(true)

	apiMux := mux.PathPrefix("/api/").Subrouter()
	apiMux.Handle("/subscription", chain.Then(subscription.Handler()))

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

func debug(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%q", dump)
		h.ServeHTTP(w, r)
	})
}
