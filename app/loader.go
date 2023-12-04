package app

import (
	"scrm-sync-data/app/config"
	"scrm-sync-data/app/core"
)

type Loader struct {
}

// LoadConfig
// @Description: 从配置文件中加载整体的应用配置
// @receiver loader
// @return error
func (loader *Loader) LoadConfig(config *config.Config) error {
	return config.Load()
}

// SaveConfig
// @Description: 从配置文件中加载整体的应用配置
// @receiver loader
// @return error
func (loader *Loader) SaveConfig(config *config.Config) error {
	return config.Save()
}

// LoadCore
// @Description: 加载整体应用核心
// @receiver loader
func (loader *Loader) LoadCore(config *config.Config, core *core.Core) {
	core.Init(config)
}
