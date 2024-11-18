package main

import (
	"os"
	"path/filepath"
	"sifu-box/execute"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/route"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func init() {
	utils.SetValue(utils.GetProjectDir(),"project-dir")
	utils.GetCore()
	utils.GetDatabase()
	if err := utils.LoadConfig(filepath.Join("config","mode.config.yaml"),"mode"); err != nil {
		utils.LoggerCaller("加载服务模式配置失败",err,1)
		os.Exit(2)
	}
	utils.LoggerCaller("加载服务模式配置完成",nil,1)
	if err := utils.LoadConfig(filepath.Join("config","proxy.config.yaml"),"proxy"); err != nil {
		utils.LoggerCaller("加载代理集合配置失败",err,1)
		os.Exit(2)
	}
	utils.LoggerCaller("加载代理集合配置完成",nil,1)
	if err := utils.LoadTemplate(); err != nil {
		utils.LoggerCaller("加载模板配置失败",err,1)
		os.Exit(2)
	}
	utils.LoggerCaller("加载模板配置完成",nil,1)
	utils.LoggerCaller("服务启动成功",nil,1)
}
func main() {
	serverMode,err := utils.GetValue("mode")
	if err != nil {
		utils.LoggerCaller("获取服务模式失败",err,1)
		os.Exit(2)
	}
	if serverMode.(models.Server).Mode {
		var lock sync.Mutex
		cronTask := cron.New()
		cronId,_ := cronTask.AddFunc("30 4 * * 1",func() {
			utils.LoggerCaller("定时任务启动",nil,1)
			var hosts []models.Host
			if err := utils.DiskDb.Find(&hosts).Error; err != nil {
				utils.LoggerCaller("获取主机集合失败",err,1)
				return
			}
			var providers []models.Provider
			if err := utils.MemoryDb.Find(&providers).Error; err != nil {
				utils.LoggerCaller("获取代理集合失败",err,1)
				return 
			}
			execute.GroupUpdate(hosts,providers,&lock,true)
		})
		cronTask.Start()
		gin.SetMode(gin.ReleaseMode)
		server := gin.Default()
		server.Use(middleware.Logger(),middleware.Recovery(true),cors.New(middleware.Cors()))
		route.SettingPages(server)
		apiGroup := server.Group("/api")
		apiGroup.GET("verify",middleware.TokenAuth())
		route.SettingHost(apiGroup)
		route.SettingFiles(apiGroup)
		route.SettingProxy(apiGroup,&lock)
		route.SettingMigrate(apiGroup)
		route.SettingTemplates(apiGroup)
		route.SettingExec(apiGroup,&lock,cronTask,&cronId)
		server.Run(serverMode.(models.Server).Listen)
	}else{
		singbox.Workflow()
	}
}