package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/polar-bear-cu/sgt-noti-service/models"
	"github.com/polar-bear-cu/sgt-noti-service/repositories"
)

var ErrAlreadySent = errors.New("reminder already sent")

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type Email struct {
	ReminderID string
	To         string
	Title      string
	Content    string
}

type EmailNotificationUsecase struct {
	mailer Mailer
	emails repositories.EmailNotificationRepository
}

func NewEmailNotification(mailer Mailer, emails repositories.EmailNotificationRepository) *EmailNotificationUsecase {
	return &EmailNotificationUsecase{mailer: mailer, emails: emails}
}

func (u *EmailNotificationUsecase) AlreadySent(ctx context.Context, reminderID string) (bool, error) {
	if reminderID == "" {
		return false, nil
	}
	return u.emails.HasSent(ctx, reminderID)
}

// Send skips a reminder that was already accepted by the mail server, since RabbitMQ
// can redeliver and the scheduler can publish the same reminder twice in one day.
// Emails without a ReminderID are always sent.
func (u *EmailNotificationUsecase) Send(ctx context.Context, e Email) error {
	if e.ReminderID != "" {
		sent, err := u.emails.HasSent(ctx, e.ReminderID)
		if err != nil {
			return err
		}
		if sent {
			return ErrAlreadySent
		}
	}

	now := time.Now()

	sendErr := u.mailer.Send(ctx, e.To, e.Title, e.Content)

	notification := models.EmailNotification{
		ReminderID: e.ReminderID,
		To:         e.To,
		Title:      e.Title,
		CreatedAt:  now,
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
