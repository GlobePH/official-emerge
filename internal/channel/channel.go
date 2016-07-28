package channel

import (
	"io"
	"net/http"

	"github.com/jeepers-creepers/emerge/internal/sms"

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
		rc := c.redisPool.Get()
		defer rc.Close()
		psc := redis.PubSubConn{rc}
		psc.Subscribe(sms.Inbox)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				if _, err := ws.Write(v.Data); err != nil {
					break
				}
			default:
				break
			}
			io.Copy(ws, ws)
		}
	}).ServeHTTP(w, r)
}
