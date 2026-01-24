package app

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
