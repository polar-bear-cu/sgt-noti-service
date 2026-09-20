package models

import "time"

const (
	StatusSent   = "sent"
	StatusFailed = "failed"
)

type NotificationLog struct {
	ID        string
	To        string
	Title     string
	Content   string
	Status    string
	CreatedAt time.Time
	SentAt    time.Time
}
