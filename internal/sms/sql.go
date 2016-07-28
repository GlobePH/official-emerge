package sms

import (
	"github.com/jackc/pgx"
)

var (
	Inbox = "sms_inbox"
)

const (
	inboxSQL = `INSERT INTO subscriber_sms_inbox (message_id, subscriber_number, message, received_at) VALUES ($1, $2, $3, $4);`
)

func PrepareStatements(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare(Inbox, inboxSQL)
	return
}
