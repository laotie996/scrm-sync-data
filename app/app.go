package app

import (
	"os"
	"os/signal"
	"scrm-sync-data/app/config"
	"scrm-sync-data/app/console"
	"scrm-sync-data/app/core"
	"scrm-sync-data/app/core/services"
	"sync"
	"syscall"
)

var once sync.Once

type App struct {
	config  config.Config
	loader  Loader
	console console.Console
	core    core.Core
}

// LoadConfig
// @Description: 加载整体应用配置
// @receiver app
// @return *App
func (app *App) LoadConfig() *App {
	err := app.loader.LoadConfig(&app.config)
	if err != nil {
		panic(err)
	}
	return app
}

// LoadCore
// @Description: 加载整体应用核心
// @receiver app
// @return *App
func (app *App) LoadCore() *App {
	app.loader.LoadCore(&app.config, &app.core)
	return app
}

// Services
// @Description: 加载整体应用核心
// @receiver app
// @return *services.Services
func (app *App) Services() *services.Services {
	return app.core.Services
}

// Init
// @Description: 应用初始化
// @receiver app
func (app *App) Init() {
	once.Do(func() {
		app.LoadConfig()
		app.LoadCore()
	})
}

// RunConsole
// @Description: 后台任务运行
// @receiver app
func (app *App) RunConsole() {
	app.console.Inject(&app.config, &app.core)
	go app.console.SyncData()
	closeChan := make(chan os.Signal, 1)
	//signal.Notify(closeChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTSTP)
	signal.Notify(closeChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	for { //接受系统信号量，从而关闭守护进程
		select {
		case <-closeChan:
			os.Exit(0)
		}
	}
}
