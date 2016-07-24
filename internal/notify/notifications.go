package notify

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type Notifier struct {
	mx        sync.Mutex
	listeners map[string]*websocket.Conn
}

func New() *Notifier {
	return &Notifier{
		listeners: make(map[string]*websocket.Conn),
	}
}

func (n *Notifier) Subscribe(conn *websocket.Conn) string {
	n.mx.Lock()
	defer n.mx.Unlock()
	id := uuid.NewV4().String()
	n.listeners[id] = conn
	return id
}

func (n *Notifier) Unsubscribe(id string) {
	n.mx.Lock()
	defer n.mx.Unlock()
	if conn, ok := n.listeners[id]; ok {
		conn.Close()
		delete(n.listeners, id)
	}
}

func (n *Notifier) Publish(data []byte) {
	n.mx.Lock()
	defer n.mx.Unlock()
	for id, conn := range n.listeners {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Print(err)
			n.Unsubscribe(id)
		}
	}
}
