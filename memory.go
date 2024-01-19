package cache

import (
	"context"
	"runtime"
	"sync"
	"time"
)

var _ Backend = (*MemoryBackend)(nil)

type Item struct {
	Value      []byte
	Expiration int64
}

func (i Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

type MemoryBackend struct {
	mu      sync.RWMutex
	items   map[string]Item
	janitor *janitor
}

func NewMemoryBackendFrom(cleanupInterval time.Duration, items map[string]Item) *MemoryBackend {
	b := &MemoryBackend{items: items}

	if cleanupInterval > 0 {
		runJanitor(b, cleanupInterval)
		runtime.SetFinalizer(b, stopJanitor)
	}

	return b
}

func NewMemoryBackend(cleanupInterval time.Duration) *MemoryBackend {
	return NewMemoryBackendFrom(cleanupInterval, make(map[string]Item))
}

func (b *MemoryBackend) Get(_ context.Context, key string) ([]byte, error) {
	b.mu.RLock()
	i, ok := b.items[key]
	b.mu.RUnlock()

	if !ok || i.Expired() {
		return nil, ErrNotFound
	}

	return i.Value, nil
}

func (b *MemoryBackend) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	var e int64
	if ttl > 0 {
		e = time.Now().Add(ttl).UnixNano()
	}

	b.mu.Lock()
	b.items[key] = Item{Value: value, Expiration: e}
	b.mu.Unlock()

	return nil
}

func (b *MemoryBackend) Del(_ context.Context, key string) error {
	b.mu.Lock()
	delete(b.items, key)
	b.mu.Unlock()

	return nil
}

func (b *MemoryBackend) DelAll(context.Context) error {
	b.mu.Lock()
	b.items = map[string]Item{}
	b.mu.Unlock()

	return nil
}

func (b *MemoryBackend) DelExpired() {
	now := time.Now().UnixNano()

	b.mu.Lock()
	for k, v := range b.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(b.items, k)
		}
	}
	b.mu.Unlock()
}

func (b *MemoryBackend) Items() map[string]Item {
	m := make(map[string]Item, len(b.items))
	now := time.Now().UnixNano()

	b.mu.RLock()
	for k, v := range b.items {
		if v.Expiration > 0 && now > v.Expiration {
			continue
		}
		m[k] = v
	}
	b.mu.RUnlock()

	return m
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(b *MemoryBackend) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			b.DelExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(b *MemoryBackend) {
	b.janitor.stop <- true
}

func runJanitor(b *MemoryBackend, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	b.janitor = j
	go j.Run(b)
}
