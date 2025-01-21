package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/JIeeiroSst/hex/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) ports.CacheRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // set if required
		DB:       0,
	})

	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *RedisCache) Get(key string) ([]byte, error) {
	return c.client.Get(c.ctx, key).Bytes()
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, data, expiration).Err()
}

func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}
