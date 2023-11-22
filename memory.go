package cache

import (
	"context"
	"sync"
	"time"
)

var _ Backend = (*MemoryBackend)(nil)

type item struct {
	value []byte
	date  time.Time
}

type MemoryBackend struct {
	mu   sync.RWMutex
	data map[string]item
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{data: make(map[string]item)}
}

func (b *MemoryBackend) Get(ctx context.Context, key string) ([]byte, error) {
	b.mu.RLock()
	i, ok := b.data[key]
	b.mu.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	if i.date.Before(time.Now()) {
		_ = b.Del(ctx, key)
		return nil, ErrExpired
	}

	return i.value, nil
}

func (b *MemoryBackend) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data[key] = item{value: value, date: time.Now().Add(ttl)}

	return nil
}

func (b *MemoryBackend) Del(_ context.Context, key string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.data, key)

	return nil
}

func (b *MemoryBackend) DelAll(context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	clear(b.data)

	return nil
}

func (b *MemoryBackend) Close() error {
	return nil
}
