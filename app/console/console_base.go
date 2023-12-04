package console

import (
	"scrm-sync-data/app/config"
	"scrm-sync-data/app/core"
	"sync"
)

type Console struct {
	SyncDataConsole
}

var this struct {
	*core.Core
}

var once sync.Once

func (Console Console) Inject(config *config.Config, core *core.Core) {
	once.Do(func() {
		this.Core = core.Clone()
	})
}
