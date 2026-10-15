package consumer

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/polar-bear-cu/sgt-noti-service/dtos"
	"github.com/polar-bear-cu/sgt-noti-service/usecases"
)

const UserDeletedQueue = "user_deleted"

type UserDeletedConsumer struct {
	conn  *amqp.Connection
	email *usecases.EmailNotificationUsecase
	inApp *usecases.InAppNotificationUsecase
}

func NewUserDeletedConsumer(
	conn *amqp.Connection,
	email *usecases.EmailNotificationUsecase,
	inApp *usecases.InAppNotificationUsecase,
) *UserDeletedConsumer {
	return &UserDeletedConsumer{conn: conn, email: email, inApp: inApp}
}

func (c *UserDeletedConsumer) Run(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	q, err := ch.QueueDeclare(UserDeletedQueue, true, false, false, false, nil)
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

func (c *UserDeletedConsumer) handle(ctx context.Context, d amqp.Delivery) {
	var msg dtos.UserDeletedMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		log.Printf("user_deleted: bad message: %v", err)
		_ = d.Nack(false, false)
		return
	}
	if msg.UserID == "" || msg.Email == "" {
		log.Printf("user_deleted: bad message: user_id and email are required")
		_ = d.Nack(false, false)
		return
	}

	inApp, err := c.inApp.DeleteAll(ctx, msg.UserID)
	if err != nil {
		log.Printf("user_deleted: delete in-app failed user_id=%s: %v", msg.UserID, err)
		_ = d.Nack(false, true)
		return
	}

	emails, err := c.email.DeleteByRecipient(ctx, msg.Email)
	if err != nil {
		log.Printf("user_deleted: delete email logs failed user_id=%s: %v", msg.UserID, err)
		_ = d.Nack(false, true)
		return
	}

	log.Printf("user_deleted: removed %d in-app, %d email logs user_id=%s", inApp, emails, msg.UserID)
	_ = d.Ack(false)
}
