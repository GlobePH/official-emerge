package subscription

import (
	"github.com/jackc/pgx"
)

type subscribers struct {
	pool *pgx.ConnPool
}

func NewSubscribers(pool *pgx.ConnPool) *subscribers {
	return &subscribers{
		pool: pool,
	}
}

func (ss *subscribers) Add(s Subscriber) (err error) {
	_, err = ss.pool.Exec(Subscribe, s.SubscriberNumber, s.AccessToken)
	return
}

func (ss *subscribers) Get(subscriberNumber string) (s *Subscriber, err error) {
	//cmd := "SELECT subscriber_number, access_token, subscription_start, subscription_end FROM subscribers WHERE subscriber_number = $1;"
	return
}
