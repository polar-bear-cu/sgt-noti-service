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

func (u *NotificationUsecase) Send(ctx context.Context, to, subject, body string) error {
	sendErr := u.mailer.Send(ctx, to, subject, body)

	log := models.NotificationLog{
		To:      to,
		Subject: subject,
		Body:    body,
		Status:  models.StatusSent,
		SentAt:  time.Now(),
	}
	if sendErr != nil {
		log.Status = models.StatusFailed
		log.Error = sendErr.Error()
	}

	if err := u.logs.Create(ctx, log); err != nil {
		return err
	}
	return sendErr
}
