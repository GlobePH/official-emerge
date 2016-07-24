package notify

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

func Handler(db *sql.DB) http.Handler {
	return handlers.MethodHandler{
		http.MethodGet:  get(db),
		http.MethodPost: post(db),
	}
}

func get(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
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

		const layout string = "Mon Jan 2 2006 15:04:05 MST+0000 (UTC)"
		for _, msg := range data.InboundSMSMessageList.InboundSMSMessage {
			receivedAt, err := time.Parse(layout, msg.DateTime)
			if err != nil {
				panic(err)
			}

			if err := saveSMSMessage(msg.MessageId, msg.SenderAddress[:10], msg.Message, receivedAt, db); err != nil {
				panic(err)
			}
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			panic(err)
		}
	})
}

func saveSMSMessage(messageId, subscriberNumber, message string, receivedAt time.Time, db *sql.DB) (err error) {
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

	if _, err = stmt.Exec(messageId, subscriberNumber, message, receivedAt); err != nil {
		return
	}

	if err = stmt.Close(); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
	return nil
}
