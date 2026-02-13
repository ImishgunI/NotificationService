package store

import (
	"context"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisStore struct {
	client redis.Client
}

func GetUrlString() string {
	sb := strings.Builder{}
	sb.Grow(1)
	sb.WriteString("redis://:" + viper.GetString("REDIS_PASSWORD") + "@" + viper.GetString("REDIS_HOST") + ":" + viper.GetString("REDIS_PORT") + "/" + viper.GetString("REDIS_DB_NUM"))
	return sb.String()
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

func (r *RedisStore) Close() {
	if err := r.client.Close(); err != nil {
		log.Fatal(err)
	}
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
