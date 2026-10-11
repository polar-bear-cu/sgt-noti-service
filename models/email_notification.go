package models

import "time"

const (
	StatusSent   = "sent"
	StatusFailed = "failed"
)

type EmailNotification struct {
	To        string
	Title     string
	Content   string
	Status    string
	CreatedAt time.Time
	SentAt    *time.Time
}
