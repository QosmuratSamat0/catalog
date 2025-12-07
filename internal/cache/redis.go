package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Cache{client: client}
}

func (c *Cache) Set(ctx context.Context, key string, value string, tll time.Duration) error {
	return c.client.Set(ctx, key, value, tll).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}
