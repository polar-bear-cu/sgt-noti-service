package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/polar-bear-cu/sgt-noti-service/models"
)

type EmailNotificationRepository interface {
	Create(ctx context.Context, notification models.EmailNotification) error
}

type EmailNotificationMongo struct {
	col *mongo.Collection
}

func NewEmailNotificationMongo(db *mongo.Database) *EmailNotificationMongo {
	return &EmailNotificationMongo{col: db.Collection("email_logs")}
}

func (r *EmailNotificationMongo) Create(ctx context.Context, notification models.EmailNotification) error {
	_, err := r.col.InsertOne(ctx, bson.M{
		"to":         notification.To,
		"title":      notification.Title,
		"content":    notification.Content,
		"status":     notification.Status,
		"created_at": notification.CreatedAt,
		"sent_at":    notification.SentAt,
	})
	return err
}
