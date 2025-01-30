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
	"sifu-box/utils"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

var taskLogger *zap.Logger
var buntClient *buntdb.DB
var entClient *ent.Client
var setting *models.Setting

var scheduler *gocron.Scheduler
func init() {
	var err error
	cmd.Execute()
	initLogger := initial.GetLogger(cmd.WorkDir, "init")
	defer initLogger.Sync()
	taskLogger = initial.GetLogger(cmd.WorkDir, "task")

	buntClient = initial.InitBuntdb(initLogger)
	initLogger.Info("内存数据库BuntDB初始化完成")

	setting, err = initial.InitSetting(cmd.Config, cmd.Server, buntClient, initLogger)
	if err != nil {
		panic(err)
	}
	initLogger.Info("Singbox配置初始化完成")

	if cmd.Server {
		entClient = initial.InitEntdb(cmd.WorkDir, initLogger)

		initLogger.Info("加载配置文件完成")
		if err := initial.SetDefaultTemplate(cmd.WorkDir, buntClient, initLogger); err != nil {
			panic(err)
		}
		initLogger.Info("读取默认模板完成")
		scheduler = gocron.NewScheduler(time.Local)
		scheduler.StartAsync()
		initLogger.Info("定时器初始化完成")
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
	if cmd.Server { webLogger = initial.GetLogger(cmd.WorkDir, "web") }
	defer func() {
		taskLogger.Sync()
		buntClient.Close()
		if cmd.Server { webLogger.Sync() }
		if entClient != nil {entClient.Close()}
	}()

	if cmd.Server {
		utils.SetValue(buntClient, models.CURRENTPROVIDER, "夜煞云", taskLogger)
		utils.SetValue(buntClient, models.CURRENTTEMPLATE, "default", taskLogger)
		singbox.GenerateConfigFiles(entClient, buntClient, nil, nil, cmd.WorkDir, cmd.Server, &rwLock, taskLogger)
		if _, err := scheduler.Cron(setting.Application.Server.Interval).Do(func(){
			singbox.GenerateConfigFiles(entClient, buntClient, nil, nil, cmd.WorkDir, cmd.Server, &rwLock, taskLogger)
			singbox.ApplyNewConfig(cmd.WorkDir, *setting.Application.Singbox, buntClient, &rwLock, &execLock, taskLogger)
		}); err != nil {
			taskLogger.Error(fmt.Sprintf("设置定时任务失败: [%s]", err.Error()))
			panic(err)
		}
		gin.SetMode(gin.ReleaseMode)
		server := gin.Default()
		server.Use(middleware.Logger(webLogger),middleware.Recovery(true, webLogger), cors.New(middleware.Cors()))
		api := server.Group("/api")
		route.SettingExec(api, entClient, buntClient, cmd.WorkDir, setting.Application.Server.User, &execLock, &rwLock, setting.Application.Singbox, webLogger)
		route.SettingFiles(api, setting.Application.Server.User, cmd.WorkDir, entClient, webLogger)
		route.SettingLogin(api, setting.Application.Server.User, webLogger)
		route.SettingConfiguration(api, cmd.WorkDir, entClient, *setting.Application.Server.User, buntClient, &rwLock, &execLock, *setting.Application.Singbox, webLogger)
		if setting.Application.Server.SSL != nil {
			fmt.Println(setting.Application.Server.SSL)
			server.Run(cmd.Listen)
		}else{
			server.Run(cmd.Listen)
		}
	}else{
		singbox.GenerateConfigFiles(nil, buntClient, nil, nil, cmd.WorkDir, cmd.Server, &rwLock, taskLogger)
	}
	
	
	
	
}

// func getWorkDir() (string, error) {
// 	// workDir, err := os.Executable()
	
// 	// workDir := "E:/Myproject/sifu-box@1.1.0/bin"
// 	var err error
// 	workDir := "/opt/sifubox/bin/bin"
// 	return filepath.Dir(filepath.Dir(workDir)), err
// }


