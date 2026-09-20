package usecases

import (
	"context"
	"time"

	"github.com/polar-bear-cu/sgt-noti-service/models"
	"github.com/polar-bear-cu/sgt-noti-service/repositories"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type NotificationUsecase struct {
	mailer Mailer
	logs   repositories.NotificationLogRepository
}

func NewNotification(mailer Mailer, logs repositories.NotificationLogRepository) *NotificationUsecase {
	return &NotificationUsecase{mailer: mailer, logs: logs}
}

func (u *NotificationUsecase) Send(ctx context.Context, to, title, content string) error {
	now := time.Now()

	sendErr := u.mailer.Send(ctx, to, title, content)

	log := models.NotificationLog{
		To:        to,
		Title:     title,
		Content:   content,
		Status:    models.StatusSent,
		CreatedAt: now,
		SentAt:    now,
	}

	if sendErr != nil {
		log.Status = models.StatusFailed
	}

	if err := u.logs.Create(ctx, log); err != nil {
		return err
	}

	return sendErr
}
