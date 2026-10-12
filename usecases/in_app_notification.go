package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/polar-bear-cu/sgt-noti-service/models"
	"github.com/polar-bear-cu/sgt-noti-service/repositories"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrSkipInApp            = errors.New("skip in-app notification")
)

const (
	maxInAppTitle   = 100
	maxInAppContent = 300
)

type InAppInput struct {
	ReminderID     string
	UserID         string
	SubscriptionID string
	Title          string
	Content        string
}

type InAppNotificationUsecase struct {
	repo repositories.InAppNotificationRepository
}

func NewInAppNotification(repo repositories.InAppNotificationRepository) *InAppNotificationUsecase {
	return &InAppNotificationUsecase{repo: repo}
}

func (u *InAppNotificationUsecase) Record(ctx context.Context, in InAppInput) error {
	if err := validateInApp(in); err != nil {
		return err
	}

	now := time.Now()
	return u.repo.Upsert(ctx, models.InAppNotification{
		ReminderID:     in.ReminderID,
		UserID:         in.UserID,
		SubscriptionID: in.SubscriptionID,
		Title:          in.Title,
		Content:        in.Content,
		Status:         true,
		ReadStatus:     false,
		CreatedAt:      now,
		SentAt:         &now,
	})
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

func (u *InAppNotificationUsecase) UnreadCount(ctx context.Context, userID string) (int64, error) {
	return u.repo.CountUnreadByUser(ctx, userID)
}

func (u *InAppNotificationUsecase) MarkAllRead(ctx context.Context, userID string) (int64, error) {
	return u.repo.MarkAllReadByUser(ctx, userID)
}

func validateInApp(in InAppInput) error {
	required := []struct{ name, value string }{
		{"reminder_id", in.ReminderID},
		{"user_id", in.UserID},
		{"subscription_id", in.SubscriptionID},
		{"inapp_title", in.Title},
		{"inapp_content", in.Content},
	}
	for _, f := range required {
		if strings.TrimSpace(f.value) == "" {
			return fmt.Errorf("%w: missing %s", ErrSkipInApp, f.name)
		}
	}
	if utf8.RuneCountInString(in.Title) > maxInAppTitle || utf8.RuneCountInString(in.Content) > maxInAppContent {
		return fmt.Errorf("%w: title or content too long", ErrSkipInApp)
	}
	return nil
}
