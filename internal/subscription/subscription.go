package subscription

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/jackc/pgx"
)

var (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)

type subscription struct {
	ss *subscribers
}

func New(pool *pgx.ConnPool) *subscription {
	return &subscription{
		ss: NewSubscribers(pool),
	}
}

func (s *subscription) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	get := http.HandlerFunc(s.get)
	handlers.MethodHandler{
		http.MethodGet:  get,
		http.MethodHead: get,
	}.ServeHTTP(w, r)
}

// I don't know why this is a GET instead of a POST
func (s *subscription) get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	subscriber := Subscriber{
		AccessToken:      q.Get("access_token"),
		SubscriberNumber: q.Get("subscriber_number"),
	}

	if subscriber.AccessToken == "" || subscriber.SubscriberNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.ss.Add(subscriber); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusAccepted)
}
