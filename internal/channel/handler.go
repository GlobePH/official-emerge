package channel

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/handlers"
)

func Handler(db *sql.DB) http.Handler {
	return handlers.MethodHandler{
		http.MethodGet: get(db),
	}
}

func get(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
