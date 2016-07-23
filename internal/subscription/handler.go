package subscription

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

type subscriber struct {
	AccesToken        string
	SubscriberNumber  string
	SubscriptionStart time.Time
	SubscriptionEnd   time.Time
}

func Handler(db *sql.DB) http.Handler {
	return handlers.MethodHandler{
		http.MethodGet:  get(db),
		http.MethodPost: post(db),
	}
}

func get(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		s := subscriber{
			AccesToken:       q.Get("access_token"),
			SubscriberNumber: q.Get("subscriber_number"),
		}

		if s.AccesToken == "" {
			http.Error(w, "Missing data: access_token.\n", http.StatusBadRequest)
			return
		}

		if s.SubscriberNumber == "" {
			http.Error(w, "Missing data: subscriber_number.\n", http.StatusBadRequest)
			return
		}

		exists, err := exists(s, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if exists {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		if err := subscribe(s, db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("%s is subscribed.\n", s.SubscriberNumber)))
	})
}

func post(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func exists(s subscriber, db *sql.DB) (bool, error) {
	cmd := "SELECT 1 FROM subscribers WHERE subscriber_number = $1;"
	stmt, err := db.Prepare(cmd)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var x int
	err = stmt.QueryRow(s.SubscriberNumber).Scan(&x)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func subscribe(s subscriber, db *sql.DB) (err error) {
	cmd := "INSERT INTO subscribers (subscriber_number, access_token) VALUES ($1, $2);"

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(cmd)
	if err != nil {
		return
	}

	if _, err = stmt.Exec(s.SubscriberNumber, s.AccesToken); err != nil {
		return
	}

	if err = stmt.Close(); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}
