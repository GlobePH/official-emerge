package subscription

import (
	"encoding/json"
	"net/http"

	"github.com/bgentry/que-go"
	"github.com/gorilla/handlers"
)

var (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)

type subscription struct {
	q *que.Client
}

func New(c *que.Client) *subscription {
	return &subscription{
		q: c,
	}
}

func (s *subscription) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	get := http.HandlerFunc(s.get)
	handlers.MethodHandler{
		http.MethodGet:  get,
		http.MethodHead: get,
	}.ServeHTTP(w, r)
}

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

	args, err := json.Marshal(subscriber)
	if err != nil {
		panic(err)
	}

	err = s.q.Enqueue(&que.Job{
		Type: Subscribe,
		Args: args,
	})

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusAccepted)
}
