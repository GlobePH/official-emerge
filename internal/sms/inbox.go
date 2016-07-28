package sms

import (
	"github.com/jackc/pgx"
)

type inbox struct {
	pool *pgx.ConnPool
}

func NewInbox(pool *pgx.ConnPool) *inbox {
	return &inbox{
		pool: pool,
	}
}

func (i *inbox) Add(m Message) (err error) {
	_, err = i.pool.Exec(Inbox, m.Id, m.Number, m.Text, m.Received)
	return
}
