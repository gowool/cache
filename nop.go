package cache

import "context"

var _ Cache = (*NopCache)(nil)

type NopCache struct{}

func (*NopCache) Set(context.Context, string, any, ...string) error {
	return nil
}

func (*NopCache) Get(context.Context, string, any) error {
	return ErrNotFound
}

func (*NopCache) DelByKey(context.Context, string) error {
	return nil
}

func (*NopCache) DelByTag(context.Context, string) error {
	return nil
}

func (*NopCache) DelAll(context.Context) error {
	return nil
}
