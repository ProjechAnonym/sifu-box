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
		initLogger.Info("定时任务初始化完成")
		scheduler := gocron.NewScheduler(time.Local)
		_, err = scheduler.Cron(setting.Application.Server.Interval).Do(func(){
			singbox.Workflow(entClient, buntClient, nil, nil, cmd.WorkDir, cmd.Server, taskLogger)
			singbox.ApplyNewConfig(cmd.WorkDir, *setting.Application.Singbox, taskLogger)
		})
		if err != nil {
			taskLogger.Error(fmt.Sprintf("设置定时任务失败: [%s]", err.Error()))
			panic(err)
		}
		scheduler.StartAsync()
		if setting.Configuration == nil {
			initLogger.Debug("配置字段为空, 将直接使用数据库中配置")
			return
		}
		initial.SaveNewProxySetting(*setting.Configuration, entClient, initLogger)
	}

}

func main() {
	var webLogger *zap.Logger
	if cmd.Server { webLogger = initial.GetLogger(cmd.WorkDir, "web") }
	defer func() {
		taskLogger.Sync()
		buntClient.Close()
		if cmd.Server { webLogger.Sync() }
		if entClient != nil {entClient.Close()}
	}()

	if cmd.Server {
		gin.SetMode(gin.ReleaseMode)
		server := gin.Default()
		server.Use(middleware.Logger(webLogger),middleware.Recovery(true, webLogger), cors.New(middleware.Cors()))
		api := server.Group("/api")
		route.SettingLogin(api, setting.Application.Server.User, webLogger)
		route.SettingConfiguration(api, entClient, *setting.Application.Server.User, webLogger)
		if setting.Application.Server.SSL != nil {
			fmt.Println(setting.Application.Server.SSL)
			server.Run(cmd.Listen)
		}else{
			server.Run(cmd.Listen)
		}
	}else{
		singbox.Workflow(nil, buntClient, nil, nil, cmd.WorkDir, cmd.Server, taskLogger)
	}
	
	
	
	
}

// func getWorkDir() (string, error) {
// 	// workDir, err := os.Executable()
	
// 	// workDir := "E:/Myproject/sifu-box@1.1.0/bin"
// 	var err error
// 	workDir := "/opt/sifubox/bin/bin"
// 	return filepath.Dir(filepath.Dir(workDir)), err
// }


