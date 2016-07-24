package notify

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

type sms struct {
	MessageId        string
	SubscriberNumber string
	Message          string
	ReceivedAt       time.Time
}

func Handler(db *sql.DB) http.Handler {
	return handlers.MethodHandler{
		http.MethodPost: post(db),
	}
}

func post(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

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
			receivedAt, err := time.Parse(layout, msg.DateTime)
			if err != nil {
				panic(err)
			}

			smsMsg := sms{
				MessageId:        msg.MessageId,
				SubscriberNumber: msg.SenderAddress[len(msg.SenderAddress)-10:],
				Message:          msg.Message,
				ReceivedAt:       receivedAt,
			}

			log.Print(smsMsg)
			if err := saveSMS(smsMsg, db); err != nil {
				panic(err)
			}
		}
	})
}

func saveSMS(msg sms, db *sql.DB) (err error) {
	cmd := "INSERT INTO subscriber_sms_inbox (message_id, subscriber_number, message, received_at) VALUES ($1, $2, $3, $4);"

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(cmd)
	if err != nil {
		return
	}

	if _, err = stmt.Exec(msg.MessageId, msg.SubscriberNumber, msg.Message, msg.ReceivedAt); err != nil {
		return
	}

	if err = stmt.Close(); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}
