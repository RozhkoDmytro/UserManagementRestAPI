package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type CacheInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(redisURL string) *RedisClient {
	// Parse the Redis URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Could not parse Redis URL: %v", err)
	}

	// Create a new Redis client using the parsed options
	rdb := redis.NewClient(opt)

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return &RedisClient{Client: rdb}
}

// Get отримує значення за ключем з Redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Ключ не існує
		return "", nil
	} else if err != nil {
		// Помилка під час виконання запиту
		return "", err
	}
	return val, nil
}

// Set зберігає значення за ключем у Redis
func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	err := r.Client.Set(ctx, key, value, time.Minute).Err()
	if err != nil {
		// Помилка під час виконання запиту
		return err
	}
	return nil
}
