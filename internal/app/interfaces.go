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
	GetEvent(ctx context.Context, key string) (domain.Event, error)
	UpdateEvent(event *domain.Event) error 
}

type EventPublisher interface {
	PublishEvent(ctx context.Context, key string) error
}

type NotificationSender interface {
	PublishNotification()
}

type EventQueue interface {
	ConsumeEvent() (string, error)
	AckEvent()
	NackEvent()
}
