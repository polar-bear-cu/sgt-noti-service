package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/polar-bear-cu/sgt-noti-service/models"
)

type NotificationLogRepository interface {
	Create(ctx context.Context, log models.NotificationLog) error
	HasSent(ctx context.Context, reminderID string) (bool, error)
}

type NotificationLogMongo struct {
	col *mongo.Collection
}

func NewNotificationLogMongo(db *mongo.Database) *NotificationLogMongo {
	return &NotificationLogMongo{col: db.Collection("email_logs")}
}

// Check there's index of reminder_id
func (r *NotificationLogMongo) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "reminder_id", Value: 1}, {Key: "status", Value: 1}},
	})
	return err
}

func (r *NotificationLogMongo) Create(ctx context.Context, log models.NotificationLog) error {
	_, err := r.col.InsertOne(ctx, bson.M{
		"reminder_id": log.ReminderID,
		"to":          log.To,
		"title":       log.Title,
		"content":     log.Content,
		"status":      log.Status,
		"created_at":  log.CreatedAt,
		"sent_at":     log.SentAt,
	})
	return err
}

func (r *NotificationLogMongo) HasSent(ctx context.Context, reminderID string) (bool, error) {
	err := r.col.FindOne(ctx, bson.M{"reminder_id": reminderID, "status": models.StatusSent}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	return err == nil, err
}
