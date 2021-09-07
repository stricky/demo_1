package models

import (
	"time"
)

type OutgoingMessage struct {
	Time     time.Time
	Uid      int64
	BodyText string
	Hash     string
}
