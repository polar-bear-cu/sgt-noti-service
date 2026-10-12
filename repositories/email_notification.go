package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/polar-bear-cu/sgt-noti-service/models"
)

type EmailNotificationRepository interface {
	Create(ctx context.Context, notification models.EmailNotification) error
	HasSent(ctx context.Context, reminderID string) (bool, error)
}

type EmailNotificationMongo struct {
	col *mongo.Collection
}

func NewEmailNotificationMongo(db *mongo.Database) *EmailNotificationMongo {
	return &EmailNotificationMongo{col: db.Collection("email_logs")}
}

func (r *EmailNotificationMongo) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "reminder_id", Value: 1}, {Key: "status", Value: 1}},
	})
	return err
}

func (r *EmailNotificationMongo) Create(ctx context.Context, notification models.EmailNotification) error {
	_, err := r.col.InsertOne(ctx, bson.M{
		"reminder_id": notification.ReminderID,
		"to":          notification.To,
		"title":       notification.Title,
		"content":     notification.Content,
		"status":      notification.Status,
		"created_at":  notification.CreatedAt,
		"sent_at":     notification.SentAt,
	})
	return err
}

func (r *EmailNotificationMongo) HasSent(ctx context.Context, reminderID string) (bool, error) {
	err := r.col.FindOne(ctx, bson.M{"reminder_id": reminderID, "status": models.StatusSent}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	return err == nil, err
}
