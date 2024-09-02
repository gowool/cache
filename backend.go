package cache

import (
	"context"
	"time"

	"github.com/coocood/freecache"
)

var ErrNotFound = freecache.ErrNotFound

type Backend interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	DelAll(ctx context.Context) error
}
