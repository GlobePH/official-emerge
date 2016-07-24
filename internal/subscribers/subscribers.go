package subscribers

import (
	"database/sql"
	"time"
)

type Subscriber struct {
	AccesToken        string
	SubscriberNumber  string
	SubscriptionStart *time.Time
	SubscriptionEnd   *time.Time
}

type Subscribers struct {
	db *sql.DB
}

func New(db *sql.DB) *Subscribers {
	return &Subscribers{
		db: db,
	}
}

func (ss *Subscribers) Add(s Subscriber) (err error) {
	cmd := "INSERT INTO subscribers (subscriber_number, access_token) VALUES ($1, $2);"

	tx, err := ss.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(cmd)
	if err != nil {
		return
	}

	if _, err = stmt.Exec(s.SubscriberNumber, s.AccesToken); err != nil {
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

func (ss *Subscribers) Get(subscriberNumber string) (s *Subscriber, err error) {
	cmd := "SELECT subscriber_number, access_token, subscription_start, subscription_end FROM subscribers WHERE subscriber_number = $1;"
	stmt, err := ss.db.Prepare(cmd)
	if err != nil {
		return
	}
	defer stmt.Close()

	s = &Subscriber{}
	if err = stmt.QueryRow(subscriberNumber).Scan(&s.SubscriberNumber, &s.AccesToken, &s.SubscriptionStart, &s.SubscriptionEnd); err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
