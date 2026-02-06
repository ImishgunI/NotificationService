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
	_, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisStore) SaveKey(ctx context.Context, key string) error {
	err := r.client.Set(ctx, key, "1", 0).Err()
	if err != nil {
		return err
	}
	return nil
}
