package subscription

import (
	"time"
)

type Subscriber struct {
	AccessToken       string     `json:"access_token,omitempty"`
	SubscriberNumber  string     `json:"subscriber_number,omitempty"`
	SubscriptionStart *time.Time `json:"subscription_start,omitempty"`
	SubscriptionEnd   *time.Time `json:"subscription_end,omitempty"`
}
