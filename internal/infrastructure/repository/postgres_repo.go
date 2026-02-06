package repository

import (
	"NotificationService/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewPoolPG(ctx context.Context, url string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Repository{
		pool: pool,
	}, nil
}

func (p *Repository) CloseDB() {
	p.pool.Close()
}

func (p *Repository) SaveEvent(ctx context.Context, e *domain.Event) error {
	_, err := p.pool.Exec(ctx,
		`INSERT INTO events (event_key, status, payload) VALUES ($1, $2, $3)`,
		e.GetKey(), e.GetStatus(), e.GetPayload())
	if err != nil {
		return err
	}
	return nil
}

func (p *Repository) GetEvent(ctx context.Context, key string) (domain.Event, error) {
	e := domain.Event{}
	var (
		status  domain.EventStatus
		payload any
	)
	err := p.pool.QueryRow(ctx, `
		SELECT status, payload FROM events
		WHERE event_key = $1
	`, key).Scan(&status, &payload)
	if err != nil {
		return domain.Event{}, err
	}
	e.SetKey(key)
	e.SetStatus(status)
	e.SetPayload(payload)
	return e, nil
}

func (p *Repository) UpdateEventStatus(ctx context.Context, eventStatus domain.EventStatus, key string) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE events
		SET status = $1
		WHERE event_key = $2
	`, eventStatus, key)
	if err != nil {
		return err
	}
	return nil
}
