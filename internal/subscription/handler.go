package subscription

import (
	"encoding/json"
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
		http.MethodPost: post(),
	}
}

func post() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("No Content.\n"))
			return
		}

		s := subscription{}
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("JSON error: %v\n", err)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s successfully subscribed.\n", s.SubscriberNumber)))
	})
}
