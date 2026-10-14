package models

import "time"

const (
	StatusSent   = "sent"
	StatusFailed = "failed"
)

type EmailNotification struct {
	ReminderID string
	To         string
	Title      string
	Status     string
	CreatedAt  time.Time
	SentAt     *time.Time
}
