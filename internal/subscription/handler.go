package subscription

import (
	"fmt"
	"net/http"

	"github.com/jeepers-creepers/emerge/internal/subscribers"

	"github.com/gorilla/handlers"
)

func Handler(ss *subscribers.Subscribers) http.Handler {
	return handlers.MethodHandler{
		http.MethodGet:  get(ss),
		http.MethodPost: post(ss),
	}
}

func get(ss *subscribers.Subscribers) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		subscriberNumber := q.Get("subscriber_number")
		if subscriberNumber == "" {
			http.Error(w, "Missing data: subscriber_number.\n", http.StatusBadRequest)
			return
		}

		s, err := ss.Get(subscriberNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if s != nil {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		s = &subscribers.Subscriber{
			SubscriberNumber: subscriberNumber,
			AccesToken:       q.Get("access_token"),
		}
		if s.AccesToken == "" {
			http.Error(w, "Missing data: access_token.\n", http.StatusBadRequest)
			return
		}

		if err := ss.Add(*s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("%s is subscribed.\n", s.SubscriberNumber)))
	})
}

func post(ss *subscribers.Subscribers) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
