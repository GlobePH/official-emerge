package notify

type Notifier struct {
	listeners  map[*listener]bool
	broadcast  chan []byte
	register   chan *listener
	unregister chan *listener
}

func New() *Notifier {
	return &Notifier{
		listeners:  make(map[*listener]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *listener),
		unregister: make(chan *listener),
	}
}

func (n *Notifier) Run() {
	for {
		select {
		case listener := <-n.register:
			n.listeners[listener] = true
		case listener := <-n.unregister:
			if _, ok := n.listeners[listener]; ok {
				delete(n.listeners, listener)
				close(listener.send)
			}
		case message := <-n.broadcast:
			for listener := range n.listeners {
				select {
				case listener.send <- message:
				default:
					close(listener.send)
					delete(n.listeners, listener)
				}
			}
		}
	}
}

func (n *Notifier) Publish(data []byte) {
	n.broadcast <- data
}
