package subscription

import (
	"github.com/jackc/pgx"
)

const (
	subscribeSQL = `INSERT INTO subscribers (subscriber_number, access_token) VALUES ($1, $2);`
)

func PrepareStatements(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare(Subscribe, subscribeSQL)
	return
}
