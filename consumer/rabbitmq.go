package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/polar-bear-cu/sgt-noti-service/dtos"
	"github.com/polar-bear-cu/sgt-noti-service/usecases"
)

// Same as Queue Name in scheduler publisher
const QueueName = "email_notifications"

type RabbitMQConsumer struct {
	conn *amqp.Connection
	uc   *usecases.NotificationUsecase
}

func NewRabbitMQConsumer(conn *amqp.Connection, uc *usecases.NotificationUsecase) *RabbitMQConsumer {
	return &RabbitMQConsumer{conn: conn, uc: uc}
}

func (c *RabbitMQConsumer) Run(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	q, err := ch.QueueDeclare(QueueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return nil
			}
			c.handle(ctx, d)
		}
	}
}

func (c *RabbitMQConsumer) handle(ctx context.Context, d amqp.Delivery) {
	var msg dtos.EmailMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		log.Printf("email_notifications: bad message: %v", err)
		_ = d.Nack(false, false)
		return
	}

	err := c.uc.Send(ctx, usecases.Email{
		ReminderID: msg.ReminderID,
		To:         msg.To,
		Title:      msg.Title,
		Content:    msg.Content,
	})
	if errors.Is(err, usecases.ErrAlreadySent) {
		log.Printf("email_notifications: skip already-sent reminder_id=%s", msg.ReminderID)
		_ = d.Ack(false)
		return
	}
	if err != nil {
		log.Printf("email_notifications: send failed: %v", err)
		_ = d.Nack(false, true)
		return
	}

	_ = d.Ack(false)
}
