package core

import (
	"scrm-sync-data/app/config"
	"scrm-sync-data/app/core/services"
)

type Core struct {
	Config *config.Config
	*services.Services
}

func (core *Core) Init(config *config.Config) {
	core.Config = config
	core.Services = new(services.Services)
	core.Services.Init(core.Config)
}

func (core *Core) Clone() *Core {
	cloneCore := *core
	return &cloneCore
}
