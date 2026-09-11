package models

import "time"

const (
	StatusSent   = "sent"
	StatusFailed = "failed"
)

type NotificationLog struct {
	ID      string
	To      string
	Subject string
	Body    string
	Status  string
	Error   string
	SentAt  time.Time
}
