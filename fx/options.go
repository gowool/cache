package fx

import (
	"go.uber.org/fx"

	"github.com/gowool/cache"
)

var (
	OptionCache = fx.Provide(func(cfg Config, backend cache.Backend) cache.Cache {
		cfg.setDefaults()

		return cache.NewCache(backend, cfg.ItemTTL)
	})
	OptionFreeBackend = fx.Provide(func(cfg Config) cache.Backend {
		cfg.setDefaults()

		return cache.NewFreeBackend(cfg.Size)
	})
	OptionChainBackend = fx.Provide(
		fx.Annotate(
			cache.NewChainBackend,
			fx.ParamTags(`group:"cache-backend"`),
			fx.As(new(cache.Backend)),
		),
	)
)
