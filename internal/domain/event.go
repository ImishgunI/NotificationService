package domain

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
	payload any
}

func NewEvent(key string, payload any) *Event {
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

func (e *Event) SetPayload(p any) {
	e.payload = p
}

func (e *Event) GetStatus() EventStatus {
	return e.status
}

func (e *Event) GetPayload() any {
	return e.payload
}

func (e *Event) GetKey() string {
	return e.key
}
