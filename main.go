package main

import (
	"fmt"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/route"
	"sifu-box/singbox"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

var taskLogger *zap.Logger
var buntClient *buntdb.DB
var entClient *ent.Client
var setting *models.Setting
var workDir, config, listen *string
var server *bool
func init() {
	var err error
	workDir, config, listen, server = cmd.InitCmd()
	initLogger := initial.GetLogger(*workDir, "init")
	defer initLogger.Sync()
	taskLogger = initial.GetLogger(*workDir, "task")

	buntClient = initial.InitBuntdb(initLogger)
	initLogger.Info("内存数据库BuntDB初始化完成")

	setting, err = initial.InitSetting(*config, *server, buntClient, initLogger)
	if err != nil {
		panic(err)
	}
	initLogger.Info("Singbox配置初始化完成")

	if *server {
		entClient = initial.InitEntdb(*workDir, initLogger)

		initLogger.Info("加载配置文件完成")
		if err := initial.SetDefaultTemplate(*workDir, buntClient, initLogger); err != nil {
			panic(err)
		}
		initLogger.Info("读取默认模板完成")

		if setting.Configuration == nil {
			initLogger.Debug("配置字段为空, 将直接使用数据库中配置")
			return
		}
		initial.SaveNewProxySetting(*setting.Configuration, entClient, initLogger)
	}

}

func main() {
	var webLogger *zap.Logger
	rwLock := sync.RWMutex{}
	execLock := sync.Mutex{}
	if *server { webLogger = initial.GetLogger(*workDir, "web") }
	defer func() {
		taskLogger.Sync()
		buntClient.Close()
		if *server { webLogger.Sync() }
		if entClient != nil {entClient.Close()}
	}()

	if *server {
		scheduler := cron.New()
		scheduler.Start()
		initial.SetDefautlApplication(entClient, buntClient, taskLogger)
		jobID, err := scheduler.AddFunc("30 4 * * 1", func(){
			singbox.GenerateConfigFiles(entClient, buntClient, nil, nil, *workDir, *server, &rwLock, taskLogger)
			singbox.ApplyNewConfig(*workDir, *setting.Application.Singbox, buntClient, &rwLock, &execLock, taskLogger)
		})
		if err != nil {
			taskLogger.Error(fmt.Sprintf("设置定时任务失败: [%s]", err.Error()))
			panic(err)
		}
		gin.SetMode(gin.ReleaseMode)
		server := gin.Default()
		server.Use(middleware.Logger(webLogger),middleware.Recovery(true, webLogger), cors.New(middleware.Cors()))
		route.SettingPages(server, *workDir)
		api := server.Group("/api")
		route.SettingMigrate(api, setting.Application.Server.User.PrivateKey, *workDir, *setting.Application.Singbox, &rwLock, &execLock, entClient, buntClient, scheduler, &jobID, webLogger)
		route.SettingHost(api, setting.Application.Server.User, entClient, buntClient, *setting.Application.Singbox, *workDir, &rwLock, &execLock, scheduler, &jobID, webLogger)
		route.SettingExec(api, entClient, buntClient, *workDir, setting.Application.Server.User, &execLock, &rwLock, setting.Application.Singbox, webLogger)
		route.SettingFiles(api, setting.Application.Server.User, *workDir, entClient, webLogger)
		route.SettingLogin(api, setting.Application.Server.User, webLogger)
		route.SettingConfiguration(api, *workDir, entClient, *setting.Application.Server.User, buntClient, &rwLock, &execLock, *setting.Application.Singbox, webLogger)
		if setting.Application.Server.SSL != nil {
			server.RunTLS(*listen, setting.Application.Server.SSL.Public, setting.Application.Server.SSL.Private)
		}else{
			server.Run(*listen)
		}
	}else{
		singbox.GenerateConfigFiles(nil, buntClient, nil, nil, *workDir, *server, &rwLock, taskLogger)
	}
	
	
	
	
}

// func getWorkDir() (string, error) {
// 	// workDir, err := os.Executable()
	
// 	// workDir := "E:/Myproject/sifu-box@1.1.0/bin"
// 	var err error
// 	workDir := "/opt/sifubox/bin/bin"
// 	return filepath.Dir(filepath.Dir(workDir)), err
// }


