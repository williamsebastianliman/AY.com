package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService interface {
  Set(ctx context.Context, key, value string, ttl time.Duration) error
  Get(ctx context.Context, key string) (string, error)
  Delete(ctx context.Context, key string) error
}

type RedisCache struct {
  client *redis.Client
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
  opt, err := redis.ParseURL(redisURL)
  if err != nil {
    return nil, err
  }
  client := redis.NewClient(opt)
  return &RedisCache{client: client}, nil
}

func (r *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
  return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
  return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
  return r.client.Del(ctx, key).Err()
}