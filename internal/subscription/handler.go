package subscription

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

type subscription struct {
	AccesToken       string    `json:"access_token,omitempty"`
	SubscriberNumber string    `json:"subscriber_number,omitempty"`
	TimeStamp        time.Time `json:"time_stamp,omitempty"`
}

func Handler() http.Handler {
	return handlers.MethodHandler{
		http.MethodGet:  get(),
		http.MethodPost: post(),
	}
}

func get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		s := subscription{
			AccesToken:       q.Get("access_token"),
			SubscriberNumber: q.Get("subscriber_number"),
		}

		if s.AccesToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing data: access_token.\n"))
			return
		}

		if s.SubscriberNumber == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing data: subscriber_number.\n"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%s successfully subscribed.\n", s.SubscriberNumber)))
	})
}

func post() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
