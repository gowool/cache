package fx

import "go.uber.org/fx"

const ModuleName = "cache"

var Module = fx.Module(
	ModuleName,
	OptionFreeBackend,
	OptionCache,
)
