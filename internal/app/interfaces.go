package app

import (
	"NotificationService/internal/domain"
	"context"
)

type IdempotencyStore interface {
	CheckKey(ctx context.Context, key string) (bool, error)
	SaveKey(ctx context.Context, key string) error
}

type EventRepository interface {
	SaveEvent(ctx context.Context, e *domain.Event) error
	GetEvent(ctx context.Context, key string) error
	UpdateStatus(status domain.EventStatus)
}

type EventPublisher interface {
	PublishEvent(ctx context.Context, e *domain.Event) error
}

type NotificationSender interface {
	PublishNotification()
}

type EventQueue interface {
	ConsumeEvent()
	AckEvent()
	NackEvent()
}
