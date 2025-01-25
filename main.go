package main

import (
	"fmt"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/models"
	"sifu-box/singbox"
	"time"

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
var configuration *models.Configuration
func init() {
	var err error
	cmd.Execute()
	initLogger := initial.GetLogger(cmd.WorkDir, "init")
	defer initLogger.Sync()
	taskLogger = initial.GetLogger(cmd.WorkDir, "task")

	buntClient = initial.InitBuntdb(initLogger)
	initLogger.Info("内存数据库BuntDB初始化完成")

	configuration, err = initial.InitConfigurationSetting(cmd.Config, buntClient, initLogger)
	if err != nil {
		panic(err)
	}
	initLogger.Info("Singbox配置初始化完成")

	if cmd.Server {
		entClient = initial.InitEntdb(cmd.WorkDir, initLogger)

		setting, err = initial.InitSetting(cmd.Config, buntClient, initLogger)
		if err != nil {
			panic(err)
		}
		initLogger.Info("加载配置文件完成")
		if err := initial.SetDefaultTemplate(cmd.WorkDir, buntClient, initLogger); err != nil {
			panic(err)
		}
		initial.SaveNewProxySetting(*configuration, entClient, initLogger)
	}

}

func main() {
	defer func() {
		taskLogger.Sync()
		buntClient.Close()
		if entClient != nil {entClient.Close()}
	}()
	var err error
	if cmd.Server {
		scheduler := gocron.NewScheduler(time.Local)
		_, err = scheduler.Cron(setting.Server.Interval).Do(func(){
			singbox.Workflow(entClient, buntClient, nil, nil, cmd.WorkDir, cmd.Server, taskLogger)
			singbox.TransferConfig(cmd.WorkDir, *setting.Singbox, taskLogger)
		})
		if err != nil {
			taskLogger.Error(fmt.Sprintf("设置定时任务失败: [%s]", err.Error()))
			panic(err)
		}
		scheduler.StartAsync()
		gin.SetMode(gin.ReleaseMode)
		server := gin.Default()
		if setting.Server.SSL != nil {
			fmt.Println(setting.Server.SSL)
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


