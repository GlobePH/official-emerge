package sms

import (
	"time"
)

type Message struct {
	Id       string    `json:"message_id,omitempty"`
	Number   string    `json:"subscriber_nbr,omitempty"`
	Text     string    `json:"message_text,omitempty"`
	Received time.Time `json:"received,omitempty"`
}
