package models

import (
	"time"
)

type IncomingMessage struct {
	Timestamp time.Time
	Body      []byte
}

func NewIncomingMessage(timestamp time.Time, body []byte) *IncomingMessage {
	return &IncomingMessage{
		Timestamp: timestamp,
		Body:      body,
	}
}

