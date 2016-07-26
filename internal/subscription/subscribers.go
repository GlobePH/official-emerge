package subscription

import (
	"database/sql"
)

type subscribers struct {
	db *sql.DB
}

func NewSubscribers(db *sql.DB) *subscribers {
	return &subscribers{
		db: db,
	}
}

func (ss *subscribers) Add(s Subscriber) (err error) {
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

	if _, err = stmt.Exec(s.SubscriberNumber, s.AccessToken); err != nil {
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

func (ss *subscribers) Get(subscriberNumber string) (s *Subscriber, err error) {
	cmd := "SELECT subscriber_number, access_token, subscription_start, subscription_end FROM subscribers WHERE subscriber_number = $1;"
	stmt, err := ss.db.Prepare(cmd)
	if err != nil {
		return
	}
	defer stmt.Close()

	s = &Subscriber{}
	if err = stmt.QueryRow(subscriberNumber).Scan(&s.SubscriberNumber, &s.AccessToken, &s.SubscriptionStart, &s.SubscriptionEnd); err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
