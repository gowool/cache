package fx

import (
	"go.uber.org/fx"

	"github.com/gowool/cache"
)

const ModuleName = "cache"

var Module = fx.Module(
	ModuleName,
	fx.Provide(func(cfg Config) cache.Backend {
		return cache.NewFreeBackend(cfg.Size)
	}),
	fx.Provide(func(cfg Config, backend cache.Backend) cache.Cache {
		return cache.NewCache(backend, cfg.ItemTTL)
	}),
)