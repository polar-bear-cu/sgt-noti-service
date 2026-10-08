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

type NotificationUsecase struct {
	mailer Mailer
	logs   repositories.NotificationLogRepository
}

func NewNotification(mailer Mailer, logs repositories.NotificationLogRepository) *NotificationUsecase {
	return &NotificationUsecase{mailer: mailer, logs: logs}
}

// Send skips a reminder that was already accepted by the mail server, since RabbitMQ
// can redeliver and the scheduler can publish the same reminder twice in one day.
// Emails without a ReminderID are always sent.
func (u *NotificationUsecase) Send(ctx context.Context, e Email) error {
	if e.ReminderID != "" {
		sent, err := u.logs.HasSent(ctx, e.ReminderID)
		if err != nil {
			return err
		}
		if sent {
			return ErrAlreadySent
		}
	}

	now := time.Now()

	sendErr := u.mailer.Send(ctx, e.To, e.Title, e.Content)

	log := models.NotificationLog{
		ReminderID: e.ReminderID,
		To:         e.To,
		Title:      e.Title,
		Content:    e.Content,
		CreatedAt:  now,
	}

	if sendErr != nil {
		log.Status = models.StatusFailed
		log.SentAt = nil
	} else {
		sentAt := time.Now()
		log.Status = models.StatusSent
		log.SentAt = &sentAt
	}

	if err := u.logs.Create(ctx, log); err != nil {
		return err
	}

	return sendErr
}
