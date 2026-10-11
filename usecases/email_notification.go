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

type EmailNotificationUsecase struct {
	mailer Mailer
	emails repositories.EmailNotificationRepository
}

func NewEmailNotification(mailer Mailer, emails repositories.EmailNotificationRepository) *EmailNotificationUsecase {
	return &EmailNotificationUsecase{mailer: mailer, emails: emails}
}

func (u *EmailNotificationUsecase) Send(ctx context.Context, to, title, content string) error {
	now := time.Now()

	sendErr := u.mailer.Send(ctx, to, title, content)

	notification := models.EmailNotification{
		To:        to,
		Title:     title,
		Content:   content,
		CreatedAt: now,
	}

	if sendErr != nil {
		notification.Status = models.StatusFailed
		notification.SentAt = nil
	} else {
		sentAt := time.Now()
		notification.Status = models.StatusSent
		notification.SentAt = &sentAt
	}

	if err := u.emails.Create(ctx, notification); err != nil {
		return err
	}

	return sendErr
}
