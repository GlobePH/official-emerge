package subscription

import (
	"encoding/json"
	"net/http"

	"github.com/bgentry/que-go"
	"github.com/gorilla/handlers"
	"github.com/jackc/pgx"
)

var (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)

type subscription struct {
	qc *que.Client
}

func New(pool *pgx.ConnPool) *subscription {
	return &subscription{
		qc: que.NewClient(pool),
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

	args, err := json.Marshal(subscriber)
	if err != nil {
		panic(err)
	}

	s.qc.Enqueue(&que.Job{
		Type: Subscribe,
		Args: args,
	})

	w.WriteHeader(http.StatusAccepted)
}
