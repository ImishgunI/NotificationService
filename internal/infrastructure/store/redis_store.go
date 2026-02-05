package store

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client redis.Client
}

func NewRedisClient(url string) (*RedisStore, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &RedisStore{
		client: *redis.NewClient(opt),
	}, nil
}

func (r *RedisStore) CheckKey(ctx context.Context, key string) (bool, error) {
	ok, err := r.client.Get(ctx, key).Bool()
	if err != nil {
		return ok, err
	}
	return ok, nil
}

func (r *RedisStore) SaveKey(ctx context.Context, key string) error {
	err := r.client.Set(ctx, "", key, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
