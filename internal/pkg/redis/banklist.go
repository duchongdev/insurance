package redis

import (
	"context"
	"time"

	redislib "github.com/redis/go-redis/v9"
)

const bankListKeyPrefix = "bank_list:"

// BankListCache 按渠道码缓存华安 getBankList 原始 JSON 响应。
type BankListCache struct {
	rdb *redislib.Client
	ttl time.Duration // 0 表示不过期，直至下次刷新覆盖
}

// NewBankListCache 构造银行列表缓存；ttl 为 0 时键永不过期。
func NewBankListCache(c *Client, ttl time.Duration) *BankListCache {
	return &BankListCache{rdb: c.RDB(), ttl: ttl}
}

func (b *BankListCache) key(channelCode string) string {
	return bankListKeyPrefix + channelCode
}

// Get 读取缓存；未命中返回 nil, nil。
func (b *BankListCache) Get(ctx context.Context, channelCode string) ([]byte, error) {
	data, err := b.rdb.Get(ctx, b.key(channelCode)).Bytes()
	if err == redislib.Nil {
		return nil, nil
	}
	return data, err
}

// Set 写入华安原始响应 JSON。
func (b *BankListCache) Set(ctx context.Context, channelCode string, data []byte) error {
	if b.ttl > 0 {
		return b.rdb.Set(ctx, b.key(channelCode), data, b.ttl).Err()
	}
	return b.rdb.Set(ctx, b.key(channelCode), data, 0).Err()
}
