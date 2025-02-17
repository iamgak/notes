package models

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStruct struct {
	client *redis.Client
}

func (c *RedisStruct) getRedis(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, key).Result()
	return val, err
}

func (c *RedisStruct) setRedis(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := c.client.Set(ctx, key, value, expiration).Err()
	return err
}

func (c *RedisStruct) Publish(ctx context.Context, msg []byte) error {
	// msg := []byte("New to-do item added")
	err := c.client.Publish(ctx, "todo.notifications", msg).Err()
	return err
}

func (c *RedisStruct) Subscribe(ctx context.Context) {
	// msg := []byte("New to-do item added")
	sub := c.client.Subscribe(ctx, "todo.notifications")
	ch := sub.Channel()
	for msg := range ch {
		fmt.Println("Received message:", msg)
	}
}
