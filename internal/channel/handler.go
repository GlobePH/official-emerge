package channel

import (
	"io"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/handlers"
	"golang.org/x/net/websocket"
)

type channel struct {
	redisPool *redis.Pool
}

func New(redisPool *redis.Pool) *channel {
	return &channel{
		redisPool: redisPool,
	}
}

func (c *channel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	get := http.HandlerFunc(c.get)
	handlers.MethodHandler{
		http.MethodGet: get,
	}.ServeHTTP(w, r)
}

func (c *channel) get(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(ws *websocket.Conn) {
		io.Copy(ws, ws)
	}).ServeHTTP(w, r)
}
