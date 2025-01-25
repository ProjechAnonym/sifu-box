package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/models"
	"sifu-box/singbox"
	"sifu-box/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
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

	singboxSetting, err := loadSingboxSetting(cmd.Config, buntClient, initLogger)
	if err != nil {
		panic(err)
	}
	initLogger.Info("Singbox配置初始化完成")

	if cmd.Server {
		entClient = initial.InitEntdb(cmd.WorkDir, initLogger)

		setting, err = loadSetting(cmd.Config, buntClient, initLogger)
		if err != nil {
			panic(err)
		}
		initLogger.Info("加载配置文件完成")
		if err := setDefaultTemplate(cmd.WorkDir, buntClient, initLogger); err != nil {
			panic(err)
		}
		initial.SaveNewProvidersOrRulesets(singboxSetting.Providers, singboxSetting.Rulesets, singboxSetting.Templates, entClient, initLogger)
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
			singbox.TransferConfig(cmd.WorkDir, *setting.SingboxEnv, taskLogger)
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


func setDefaultTemplate(workDir string, buntClient *buntdb.DB, logger *zap.Logger) error {
	file, err := os.Open(filepath.Join(workDir, models.STATICDIR, models.TEMPLATEDIR, models.DEFAULTTEMPLATEPATH))
	if err != nil {
		logger.Error(fmt.Sprintf("打开默认模板文件失败: [%s]",err.Error()))
		return err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("读取默认模板文件失败: [%s]",err.Error()))
		return err
	}
	var template models.Template
	if err := yaml.Unmarshal(content, &template); err != nil {
		logger.Error(fmt.Sprintf("解析默认模板文件失败: [%s]",err.Error()))
		return err
	}
	templateByte, err := json.Marshal(template)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化默认模板文件失败: [%s]",err.Error()))
		return err
	}
	if err := utils.SetValue(buntClient, models.DEFAULTTEMPLATEKEY, string(templateByte), logger); err != nil {
		logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
		return err
	}
	return nil
}

func loadSetting(workDir string, buntClient *buntdb.DB, logger *zap.Logger) (*models.Setting, error){
	file, err := os.Open(filepath.Join(workDir, models.SIFUBOXSETTINGFILE))
	if err != nil {
		logger.Error(fmt.Sprintf("打开默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("读取默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	var setting models.Setting
	if err := yaml.Unmarshal(content, &setting); err != nil {
		logger.Error(fmt.Sprintf("解析默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	settingByte, err := json.Marshal(setting)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	if err := utils.SetValue(buntClient, models.SIFUBOXSETTINGKEY, string(settingByte), logger); err != nil {
		logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
		return nil, err
	}
	return &setting, nil
}

func loadSingboxSetting(workDir string, buntClient *buntdb.DB, logger *zap.Logger) (*models.SingboxSetting, error) {
	file, err := os.Open(filepath.Join(workDir, models.SINGBOXSETTINGFILE))
	if err != nil {
		logger.Error(fmt.Sprintf("打开默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("读取默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	var singboxSetting models.SingboxSetting
	if err := yaml.Unmarshal(content, &singboxSetting); err != nil {
		logger.Error(fmt.Sprintf("解析默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	singboxSettingByte, err := json.Marshal(singboxSetting)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化默认模板文件失败: [%s]",err.Error()))
		return nil, err
	}
	if err := utils.SetValue(buntClient, models.SINGBOXSETTINGKEY, string(singboxSettingByte), logger); err != nil {
		logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
		return nil, err
	}
	return &singboxSetting, nil
}