package handler

import (
	"NotificationService/internal/domain"
	"context"
	"encoding/json"
	"fmt"
)

type BaseEvent struct {
	Type string `json:"type"`
}

type Dispatcher struct{}

func (d *Dispatcher) Handle(ctx context.Context, payload json.RawMessage) error {
	var base BaseEvent
	if err := json.Unmarshal(payload, &base); err != nil {
		return err
	}
	switch base.Type {
	case "send_email":
		break
	case "create_order":
		break
	default:
		return fmt.Errorf("%w", domain.ErrUnknownEventType)
	}
	return nil
}
