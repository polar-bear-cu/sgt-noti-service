package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/polar-bear-cu/sgt-noti-service/models"
)

var ErrInvalidNotificationID = errors.New("invalid notification id")

type InAppNotificationRepository interface {
	ListByUser(ctx context.Context, userID string) ([]models.InAppNotification, error)
	DeleteByIDAndUser(ctx context.Context, id, userID string) (bool, error)
	DeleteAllByUser(ctx context.Context, userID string) (int64, error)
}

type InAppNotificationMongo struct {
	col *mongo.Collection
}

func NewInAppNotificationMongo(db *mongo.Database) *InAppNotificationMongo {
	return &InAppNotificationMongo{col: db.Collection("in_app_notifications")}
}

func (r *InAppNotificationMongo) ListByUser(ctx context.Context, userID string) ([]models.InAppNotification, error) {
	cursor, err := r.col.Find(
		ctx,
		bson.M{"user_id": userID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	notifications := make([]models.InAppNotification, 0)
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *InAppNotificationMongo) DeleteByIDAndUser(ctx context.Context, id, userID string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, ErrInvalidNotificationID
	}

	result, err := r.col.DeleteOne(ctx, bson.M{
		"_id":     objectID,
		"user_id": userID,
	})
	if err != nil {
		return false, err
	}

	return result.DeletedCount == 1, nil
}

func (r *InAppNotificationMongo) DeleteAllByUser(ctx context.Context, userID string) (int64, error) {
	result, err := r.col.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
