package http

import (
	"NotificationService/internal/app"

	"github.com/gofiber/fiber/v3"
)

type Request struct {
	app *app.AcceptEvent
}

func NewRequest(app *app.AcceptEvent) *Request {
	return &Request{
		app: app,
	}
}

func (r *Request) CreateEvent(c fiber.Ctx) error {
	key := c.Request().Header.Peek("Idempotency-Key")
	payload := c.Body()
	err := r.app.Execute(string(key), payload)
	if err != nil {
		c.Status(500)
		return err
	}
	c.Status(201)
	return nil
}
