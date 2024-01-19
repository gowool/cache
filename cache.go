package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, tags ...string) error
	Get(ctx context.Context, key string, value interface{}) error
	DelByKey(ctx context.Context, key string) error
	DelByTag(ctx context.Context, tag string) error
	DelAll(ctx context.Context) error
}

type cache struct {
	backend Backend
	ttl     time.Duration
}

func NewCache(backend Backend, ttl time.Duration) Cache {
	return cache{backend: backend, ttl: ttl}
}

func (c cache) Set(ctx context.Context, key string, value interface{}, tags ...string) error {
	if err := c.set(ctx, "value:"+key, value); err != nil {
		return err
	}

	if err := c.set(ctx, "tags:"+key, tags); err != nil {
		return err
	}

NEXT:
	for _, tag := range tags {
		var keys []string
		if err := c.get(ctx, "keys:"+tag, &keys); err == nil {
			for _, k := range keys {
				if k == key {
					continue NEXT
				}
			}
		}

		keys = append(keys, key)
		_ = c.set(ctx, "keys:"+tag, keys)
	}

	return nil
}

func (c cache) Get(ctx context.Context, key string, value interface{}) error {
	return c.get(ctx, "value:"+key, value)
}

func (c cache) DelByKey(ctx context.Context, key string) error {
	err := c.backend.Del(ctx, "value:"+key)

	var tags []string
	err = errors.Join(err, c.get(ctx, "tags:"+key, &tags))
	err = errors.Join(err, c.backend.Del(ctx, "tags:"+key))

	for _, tag := range tags {
		var keys []string
		if err1 := c.get(ctx, "keys:"+tag, &keys); err1 != nil {
			err = errors.Join(err, err1)
			continue
		}

		if index := slices.Index(keys, key); index > -1 {
			keys = slices.Delete(keys, index, index+1)
			err = errors.Join(err, c.set(ctx, "keys:"+tag, keys))
		}
	}

	return err
}

func (c cache) DelByTag(ctx context.Context, tag string) error {
	var keys []string
	err := c.get(ctx, "keys:"+tag, &keys)
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = errors.Join(err, c.DelByKey(ctx, key))
	}

	return errors.Join(err, c.backend.Del(ctx, "keys:"+tag))
}

func (c cache) DelAll(ctx context.Context) error {
	return c.backend.DelAll(ctx)
}

func (c cache) set(ctx context.Context, key string, value interface{}) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.backend.Set(ctx, key, raw, c.ttl)
}

func (c cache) get(ctx context.Context, key string, value interface{}) error {
	raw, err := c.backend.Get(ctx, key)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()

	return decoder.Decode(value)
}
