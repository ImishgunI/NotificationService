package domain

import "encoding/json"

type EventStatus int

const (
	NEW        EventStatus = iota
	ACCEPTED               // AcceptEvent start
	PROCESSING             // Process Event
	DONE                   // Afrer Process Event
	FAILED                 // if process event started, but not done
	REJECTED               // if accept or process event started but something went wrong with infrastracure
)

type Event struct {
	key     string
	status  EventStatus
	payload json.RawMessage
}

func NewEvent(key string, payload json.RawMessage) *Event {
	return &Event{
		key:     key,
		status:  NEW,
		payload: payload,
	}
}

func (e *Event) Accept() {
	e.status = ACCEPTED
}

func (e *Event) Reject() {
	e.status = REJECTED
}

func (e *Event) Processing() {
	e.status = PROCESSING
}

func (e *Event) Done() {
	e.status = DONE
}

func (e *Event) Failed() {
	e.status = FAILED
}

func (e *Event) SetPayload(p json.RawMessage) {
	e.payload = p
}

func (e *Event) SetKey(key string) {
	e.key = key
}

func (e *Event) SetStatus(status EventStatus) {
	e.status = status
}

func (e *Event) GetStatus() EventStatus {
	return e.status
}

func (e *Event) GetPayload() json.RawMessage {
	return e.payload
}

func (e *Event) GetKey() string {
	return e.key
}
