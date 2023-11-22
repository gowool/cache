package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gowool/cache/internal/util"
)

var _ Backend = RedisBackend{}

type RedisBackend struct {
	prefix string
	client redis.UniversalClient
}

func NewRedisBackend(prefix string, client redis.UniversalClient) RedisBackend {
	return RedisBackend{prefix: prefix, client: client}
}

func (b RedisBackend) Get(ctx context.Context, key string) ([]byte, error) {
	return b.client.Get(ctx, b.prefix+key).Bytes()
}

func (b RedisBackend) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return b.client.Set(ctx, b.prefix+key, util.BytesToString(value), ttl).Err()
}

func (b RedisBackend) Del(ctx context.Context, key string) error {
	return b.client.Del(ctx, b.prefix+key).Err()
}

func (b RedisBackend) DelAll(ctx context.Context) (err error) {
	iter := b.client.Scan(ctx, 0, b.prefix+"*", 0).Iterator()

	for iter.Next(ctx) {
		err = errors.Join(err, b.client.Del(ctx, iter.Val()).Err())
	}

	return errors.Join(err, iter.Err())
}

func (b RedisBackend) Close() error {
	return b.client.Close()
}
