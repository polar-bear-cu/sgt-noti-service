package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/polar-bear-cu/sgt-noti-service/models"
)

type NotificationLogRepository interface {
	Create(ctx context.Context, log models.NotificationLog) error
}

type NotificationLogMongo struct {
	col *mongo.Collection
}

func NewNotificationLogMongo(db *mongo.Database) *NotificationLogMongo {
	return &NotificationLogMongo{col: db.Collection("email_logs")}
}

func (r *NotificationLogMongo) Create(ctx context.Context, log models.NotificationLog) error {
	_, err := r.col.InsertOne(ctx, bson.M{
		"to":      log.To,
		"subject": log.Subject,
		"body":    log.Body,
		"status":  log.Status,
		"error":   log.Error,
		"sentAt":  log.SentAt,
	})
	return err
}
