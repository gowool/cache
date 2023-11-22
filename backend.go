package cache

import (
	"context"
	"errors"
	"io"
	"time"
)

var _ Backend = ChainBackend{}

var (
	ErrExpired  = errors.New("item expired")
	ErrNotFound = errors.New("item not found")
)

type Backend interface {
	io.Closer
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	DelAll(ctx context.Context) error
}

type ChainBackend struct {
	backends []Backend
}

func NewChainBackend(backends ...Backend) ChainBackend {
	return ChainBackend{backends: backends}
}

func (b ChainBackend) Get(ctx context.Context, key string) ([]byte, error) {
	for _, backend := range b.backends {
		if value, err := backend.Get(ctx, key); err == nil {
			return value, nil
		}
	}

	return nil, ErrNotFound
}

func (b ChainBackend) Set(ctx context.Context, key string, value []byte, ttl time.Duration) (err error) {
	for _, backend := range b.backends {
		err = errors.Join(err, backend.Set(ctx, key, value, ttl))
	}
	return
}

func (b ChainBackend) Del(ctx context.Context, key string) (err error) {
	for _, backend := range b.backends {
		err = errors.Join(err, backend.Del(ctx, key))
	}
	return
}

func (b ChainBackend) DelAll(ctx context.Context) (err error) {
	for _, backend := range b.backends {
		err = errors.Join(err, backend.DelAll(ctx))
	}
	return
}

func (b ChainBackend) Close() (err error) {
	for _, backend := range b.backends {
		err = errors.Join(err, backend.Close())
	}
	return
}
