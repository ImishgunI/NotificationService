package app

import "context"

type IdempotencyStore interface {
	CheckKey(ctx context.Context, key string) (bool, error)
	SaveKey(ctx context.Context, key string) error
}

type EventRepository interface {
	SaveEvent(ctx context.Context) error
	GetEvent(ctx context.Context) error
	UpdateStatus()
}

type EventPublisher interface {
	PublishEvent(ctx context.Context) error
}

type NotificationSender interface {
	PublishNotification() 
}

type EventQueue interface {
	ConsumeEvent()
	AckEvent()
	NackEvent()
}
