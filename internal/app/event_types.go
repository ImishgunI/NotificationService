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
	err = ae.Publisher.PublishEvent(context.Background(), event.GetKey())
	if err != nil {
		event.Reject()
		return err
	}
	return nil
}

func (pe *ProcessEvent) Execute() error {
	key, err := pe.Queue.ConsumeEvent()
	if err != nil {
		return err
	}
	event, err := pe.Repo.GetEvent(context.Background(), key)
	if err != nil {
		return err
	}
	if event.GetStatus() != domain.ACCEPTED {
		pe.Queue.AckEvent()
		return nil
	}
	event.Processing()
	err = pe.Repo.UpdateEvent(&event) 	
	if err != nil {
		return err
	}
	/*
		Обработка...
	*/
	event.Done()
	err = pe.Repo.UpdateEvent(&event)
	if err != nil {
		return err
	}
	pe.Queue.AckEvent()
	return nil
}
