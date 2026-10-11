package usecases

import (
	"context"
	"errors"

	"github.com/polar-bear-cu/sgt-noti-service/models"
	"github.com/polar-bear-cu/sgt-noti-service/repositories"
)

var ErrNotificationNotFound = errors.New("notification not found")

type InAppNotificationUsecase struct {
	repo repositories.InAppNotificationRepository
}

func NewInAppNotification(repo repositories.InAppNotificationRepository) *InAppNotificationUsecase {
	return &InAppNotificationUsecase{repo: repo}
}

func (u *InAppNotificationUsecase) List(ctx context.Context, userID string) ([]models.InAppNotification, error) {
	return u.repo.ListByUser(ctx, userID)
}

func (u *InAppNotificationUsecase) DeleteOne(ctx context.Context, userID, notificationID string) error {
	deleted, err := u.repo.DeleteByIDAndUser(ctx, notificationID, userID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotificationNotFound
	}

	return nil
}

func (u *InAppNotificationUsecase) DeleteAll(ctx context.Context, userID string) (int64, error) {
	return u.repo.DeleteAllByUser(ctx, userID)
}
