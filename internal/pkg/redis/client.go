// Package redis 封装 go-redis 客户端：连接初始化、Ping 与生命周期管理。
package redis

import (
	"context"
	"fmt"
	"time"

	redislib "github.com/redis/go-redis/v9"
)

// Client 应用级 Redis 连接，供缓存、分布式锁等场景复用。
type Client struct {
	rdb *redislib.Client
}

// New 创建 Redis 客户端并在启动时 Ping 校验连通性。
func New(addr, password string, db int) (*Client, error) {
	rdb := redislib.NewClient(&redislib.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &Client{rdb: rdb}, nil
}

// Ping 检查 Redis 是否可用。
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// RDB 返回底层 go-redis 客户端，供业务直接读写。
func (c *Client) RDB() *redislib.Client {
	return c.rdb
}

// Close 关闭连接。
func (c *Client) Close() error {
	return c.rdb.Close()
}
