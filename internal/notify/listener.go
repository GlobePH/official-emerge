package notify

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type listener struct {
	n    *Notifier
	conn *websocket.Conn
	send chan []byte
}

func Listen(n *Notifier, conn *websocket.Conn) {
	l := &listener{
		n:    n,
		conn: conn,
		send: make(chan []byte),
	}
	l.n.register <- l
	go l.writePump()
	l.readPump()
}

func (l *listener) readPump() {
	defer func() {
		l.n.unregister <- l
		l.conn.Close()
	}()
	l.conn.SetReadLimit(maxMessageSize)
	l.conn.SetReadDeadline(time.Now().Add(pongWait))
	l.conn.SetPongHandler(func(string) error {
		l.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := l.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("readPump error: %v\n", err)
			}
			break
		}
		l.n.broadcast <- msg
	}
}

func (l *listener) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		l.conn.Close()
	}()

	for {
		select {
		case message, ok := <-l.send:
			if !ok {
				l.write(websocket.CloseMessage, []byte{})
				return
			}

			l.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := l.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(l.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-l.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := l.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (l *listener) write(mt int, data []byte) error {
	l.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return l.conn.WriteMessage(mt, data)
}
