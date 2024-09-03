package fx

import "time"

type Config struct {
	Size    int           `json:"size,omitempty" yaml:"size,omitempty"`
	ItemTTL time.Duration `json:"item_ttl,omitempty" yaml:"item_ttl,omitempty"`
}

func (cfg *Config) InitDefaults() {
	if cfg.Size <= 0 {
		cfg.Size = 100 * 1024 * 1024 // 100 MB
	}
	if cfg.ItemTTL == 0 {
		cfg.ItemTTL = 24 * time.Hour
	}
}
