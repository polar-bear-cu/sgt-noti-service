package dtos

import "time"

type InAppNotificationResponse struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	SubscriptionID string     `json:"subscription_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Status         bool       `json:"status"`
	ReadStatus     bool       `json:"read_status"`
	CreatedAt      time.Time  `json:"created_at"`
	SentAt         *time.Time `json:"sent_at"`
}

type DeleteNotificationsResponse struct {
	DeletedCount int64 `json:"deleted_count"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
