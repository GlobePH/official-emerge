package channel

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/handlers"
)

func Handler(db *sql.DB) http.Handler {
	return handlers.MethodHandler{}
}
