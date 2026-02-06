package app

import (
	"context"

	"NotificationService/internal/domain"
)

type IdempotencyStore interface {
	CheckKey(ctx context.Context, key string) (bool, error)
	SaveKey(ctx context.Context, key string) error
}

type EventRepository interface {
	SaveEvent(ctx context.Context, e *domain.Event) error
	GetEvent(ctx context.Context, key string) (domain.Event, error)
	UpdateEventStatus(ctx context.Context, eventStatus domain.EventStatus, key string) error
}

type EventPublisher interface {
	PublishEvent(ctx context.Context, key string) error
}

type EventQueue interface {
	ConsumeEvent(ctx context.Context) (string, error)
	AckEvent() error
	NackEvent() error
}

type EventHandler interface {
	Handle(event *domain.Event) error
}
