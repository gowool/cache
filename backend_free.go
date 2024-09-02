package cache

import (
	"context"
	"time"

	"github.com/coocood/freecache"

	"github.com/gowool/cache/internal"
)

var _ Backend = (*FreeBackend)(nil)

type FreeBackend struct {
	cache *freecache.Cache
}

func NewFreeBackend(size int) *FreeBackend {
	return &FreeBackend{cache: freecache.NewCache(size)}
}

func (b *FreeBackend) Get(_ context.Context, key string) ([]byte, error) {
	return b.cache.Get(internal.Bytes(key))
}

func (b *FreeBackend) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	return b.cache.Set(internal.Bytes(key), value, int(ttl.Seconds()))
}

func (b *FreeBackend) Del(_ context.Context, key string) error {
	_ = b.cache.Del(internal.Bytes(key))
	return nil
}

func (b *FreeBackend) DelAll(context.Context) error {
	b.cache.Clear()
	b.cache.ResetStatistics()
	return nil
}
