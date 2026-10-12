package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type InAppNotification struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ReminderID     string        `bson:"reminder_id" json:"reminder_id"`
	UserID         string        `bson:"user_id" json:"user_id"`
	SubscriptionID string        `bson:"subscription_id" json:"subscription_id"`
	Title          string        `bson:"title" json:"title"`
	Content        string        `bson:"content" json:"content"`
	Status         bool          `bson:"status" json:"status"`
	ReadStatus     bool          `bson:"read_status" json:"read_status"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	SentAt         *time.Time    `bson:"sent_at,omitempty" json:"sent_at"`
}
