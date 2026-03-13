package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Client Redis 客户端
type Client struct {
	*redis.Client
}

// New 创建 Redis 客户端
func New(addr, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Client{Client: rdb}, nil
}
