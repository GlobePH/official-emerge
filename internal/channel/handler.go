package channel

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

func Handler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		io.Copy(ws, ws)
	})
}
