package redisclient

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Redis *redis.Client
}

func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	if addr == "" {
		return nil, errors.New("no redis addr")
	}

	redisClient := new(RedisClient)
	client := redisClient.init(addr, password, db)

	redisClient.Redis = client
	return redisClient, nil
}

func (r *RedisClient) init(addr string, password string, db int) *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:         addr,
			Password:     password,
			DB:           db,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	)
	return client
}

func (r *RedisClient) Close() error {
	return r.Redis.Close()
}
