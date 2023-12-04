package main

import "scrm-sync-data/app"

// @title scrm-task
// @version 1.0
// @description scrm-sync-data 数据同步微服务 api说明文档
// @contact.name San
// @contact.email 1525461449@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 8.134.166.200:8686
// @BasePath /api
func main() {
	var application = &app.App{}
	application.Init()
	application.RunConsole()
}
