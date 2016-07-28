package notify

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jeepers-creepers/emerge/internal/sms"

	"github.com/bgentry/que-go"
	"github.com/gorilla/handlers"
	"github.com/jackc/pgx"
)

type notify struct {
	qc *que.Client
}

func New(pool *pgx.ConnPool) *notify {
	return &notify{
		qc: que.NewClient(pool),
	}
}

func (n *notify) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	post := http.HandlerFunc(n.post)
	handlers.MethodHandler{
		http.MethodPost: post,
	}.ServeHTTP(w, r)
}

func (n *notify) post(w http.ResponseWriter, r *http.Request) {
	var data struct {
		InboundSMSMessageList struct {
			InboundSMSMessage []struct {
				DateTime      string // Globe doesn't send the de facto JSON standard time format
				MessageId     string
				SenderAddress string
				Message       string
			}
		}
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		panic(err)
	}

	const layout string = "Mon Jan 2 2006 15:04:05 MST+0000 (UTC)"
	for _, msg := range data.InboundSMSMessageList.InboundSMSMessage {
		received, err := time.Parse(layout, msg.DateTime)
		if err != nil {
			panic(err)
		}

		smsMsg := sms.Message{
			Id:       msg.MessageId,
			Number:   msg.SenderAddress[len(msg.SenderAddress)-10:],
			Text:     msg.Message,
			Received: received,
		}

		args, err := json.Marshal(smsMsg)
		if err != nil {
			panic(err)
		}

		n.qc.Enqueue(&que.Job{
			Type: sms.Inbox,
			Args: args,
		})
	}

	w.WriteHeader(http.StatusAccepted)
}
