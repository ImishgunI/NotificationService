package app

import (
	"context"
	"time"

	"NotificationService/internal/domain"
)

type AcceptEvent struct {
	IdemStore IdempotencyStore
	Repo      EventRepository
	Publisher EventPublisher
}

type ProcessEvent struct {
	Repo    EventRepository
	Queue   EventQueue
	Handler EventHandler
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
		return err
	}
	err = ae.Repo.SaveEvent(context.Background(), event)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ae.Publisher.PublishEvent(ctx, event.GetKey())
	if err != nil {
		return err
	}
	return nil
}

func (pe *ProcessEvent) Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	key, err := pe.Queue.ConsumeEvent(ctx)
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
	err = pe.Repo.UpdateEventStatus(event.GetStatus())
	if err != nil {
		pe.Queue.NackEvent()
		return err
	}
	err = pe.Handler.Handle(&event)
	if err != nil {
		switch err.(type) {
		case domain.BusinessError:
			event.Failed()
			pe.Repo.UpdateEventStatus(event.GetStatus())
			pe.Queue.AckEvent()
			return nil
		case domain.RetryableError:
			pe.Queue.NackEvent()
			return err
		case domain.InfrasractureError:
			event.Reject()
			pe.Repo.UpdateEventStatus(event.GetStatus())
			pe.Queue.AckEvent()
			return nil
		default:
			pe.Queue.NackEvent()
			return err
		}
	}
	event.Done()
	err = pe.Repo.UpdateEventStatus(event.GetStatus())
	if err != nil {
		return err
	}
	pe.Queue.AckEvent()
	return nil
}
