package channel

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/jeepers-creepers/emerge/internal/notify"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Handler(n *notify.Notifier, db *sql.DB) http.Handler {
	return handlers.MethodHandler{
		http.MethodGet: get(n, db),
	}
}

func get(n *notify.Notifier, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := n.Subscribe(ws)

		for {
			// wait until closed(?)
			_, _, err := ws.ReadMessage()
			if err != nil {
				switch err {
				case io.EOF:
					//Closed
					break
				default:
					break
				}
			}
		}

		n.Unsubscribe(id)

		ws.WriteMessage(websocket.CloseMessage, []byte{})
	})
}
