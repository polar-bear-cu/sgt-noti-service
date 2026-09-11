package repositories

import (
	"context"
	"sync"

	"github.com/polar-bear-cu/sgt-noti-service/models"
)

type NotificationLogRepository interface {
	Create(ctx context.Context, log models.NotificationLog) error
}

type NotificationLogMemory struct {
	mu   sync.Mutex
	logs []models.NotificationLog
}

func NewNotificationLogMemory() *NotificationLogMemory {
	return &NotificationLogMemory{}
}

func (r *NotificationLogMemory) Create(_ context.Context, log models.NotificationLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logs = append(r.logs, log)
	return nil
}
