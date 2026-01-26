package app

import (
	"NotificationService/internal/domain"
	"context"
)

type AcceptEvent struct {
	IdemStore IdempotencyStore
	Repo      EventRepository
	Publisher EventPublisher
}

type ProcessEvent struct {
	Repo      EventRepository
	Queue     EventQueue
	Publisher NotificationSender
}

func (ae *AcceptEvent) Execute(key string, payload any) error {
	ok, err := ae.IdemStore.CheckKey(context.Background(), key)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	event := domain.NewEvent(key, payload)
	event.Accept()
	err = ae.IdemStore.SaveKey(context.Background(), key)
	if err != nil {
		event.Reject()
		return err
	}
	err = ae.Repo.SaveEvent(context.Background(), event)
	if err != nil {
		event.Reject()
		return err
	}
	err = ae.Publisher.PublishEvent(context.Background(), event)
	if err != nil {
		event.Reject()
		return err
	}
	return nil
}
