package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sifu-box/ent"
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/models"
	"sifu-box/utils"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	"github.com/natefinch/lumberjack"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var workflowLogger *zap.Logger
var buntClient *buntdb.DB
var entClient *ent.Client
var workDir string
func init() {
	var err error
	workDir = getWorkDir()
	initLogger := getLogger(workDir, "init")
	defer initLogger.Sync()
	workflowLogger = getLogger(workDir, "workflow")
	buntClient, err = buntdb.Open(":memory:")
	if err != nil {
		workflowLogger.Error(fmt.Sprintf("连接Buntdb数据库失败: [%s]",err.Error()))
		panic(err)
	}
	initLogger.Info("内存数据库BuntDB初始化完成")
	entClient, err = ent.Open(dialect.SQLite, fmt.Sprintf("file:%s/sifu-box.db?cache=shared&_fk=1", workDir))
	if err != nil {
		workflowLogger.Error(fmt.Sprintf("连接Ent数据库失败: [%s]",err.Error()))
		panic(err)
	}
	initLogger.Info("连接Ent数据库完成")
	if err = entClient.Schema.Create(context.Background()); err != nil {
		workflowLogger.Error(fmt.Sprintf("创建表资源失败: [%s]",err.Error()))
		panic(err)
	}
	initLogger.Info("自动迁移Ent数据库完成")
	if err := loadSetting(workDir, buntClient, initLogger); err != nil {
		panic(err)
	}
	initLogger.Info("加载配置文件完成")
	if err := setDefaultTemplate(workDir, buntClient, initLogger); err != nil {
		panic(err)
	}
	settingStr, err := utils.GetValue(buntClient, "setting", initLogger)
	if err != nil {
		initLogger.Error(fmt.Sprintf("获取配置文件失败: [%s]",err.Error()))
		panic(err)
	}
	var setting models.Setting
	if err := json.Unmarshal([]byte(settingStr), &setting); err != nil {
		initLogger.Error(fmt.Sprintf("解析配置文件失败: [%s]",err.Error()))
		panic(err)
	}
	if setting.Server.Enabled {
		for _, supplier := range setting.Providers {
			exist, err := entClient.Provider.Query().Where(provider.NameEQ(supplier.Name)).Exist(context.Background())
			if err != nil {
				initLogger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
			}
			if !exist {
				if _, err := entClient.Provider.Create().SetName(supplier.Name).SetDetour(supplier.Detour).SetPath(supplier.Path).SetRemote(supplier.Remote).Save(context.Background()); err != nil {
					initLogger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
				}
			}	
		}
		initLogger.Info("数据库写入机场信息完成")
		for _, ruleCollection:= range setting.Rulesets {
			exist, err := entClient.RuleSet.Query().Where(ruleset.TagEQ(ruleCollection.Tag)).Exist(context.Background())
			if err != nil {
				initLogger.Error(fmt.Sprintf("获取数据库数据失败: [%s]",err.Error()))
			}
			if !exist {
				if _, err := entClient.RuleSet.Create().SetTag(ruleCollection.Tag).
														SetOutbound(ruleCollection.Outbound).
														SetPath(ruleCollection.Path).
														SetType(ruleCollection.Type).
														SetFormat(ruleCollection.Format).
														SetChina(ruleCollection.China).
														SetLabel(ruleCollection.Label).
														SetDownloadDetour(ruleCollection.DownloadDetour).
														SetUpdateInterval(ruleCollection.UpdateInterval).
														Save(context.Background()); err != nil {
					initLogger.Error(fmt.Sprintf("保存数据失败: [%s]", err.Error()))
				}
			}
		}
		initLogger.Info("数据库写入规则集信息完成")
	}
}

func main() {
	defer func() {
		workflowLogger.Sync()
		buntClient.Close()
		entClient.Close()
	}()

}

func getWorkDir() string {
	// workDir := filepath.Dir(os.Args[0])
	workDir := "E:/Myproject/sifu-box@1.1.0/bin"
	return filepath.Dir(workDir)
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriter(level, task, workDir string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/logs/sifu-box-%s-%s.log", workDir, task, level),
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLogger(workDir, task string) *zap.Logger{
	encoder := getEncoder()
	infoWriter := getWriter("info", task, workDir)
	errorWriter := getWriter("error", task, workDir)
	infoCore := zapcore.NewCore(encoder, infoWriter, zapcore.InfoLevel)
	errorCore := zapcore.NewCore(encoder, errorWriter, zapcore.ErrorLevel)
	consoloCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	core := zapcore.NewTee(infoCore, errorCore, consoloCore)
	logger := zap.New(core,zap.AddCaller())
	return logger
}

func setDefaultTemplate(workDir string, buntClient *buntdb.DB, logger *zap.Logger) error {
	file, err := os.Open(filepath.Join(workDir, "static", "default.template.yaml"))
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
	if err := utils.SetValue(buntClient, "template:default", string(templateByte), logger); err != nil {
		logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
		return err
	}
	return nil
}

func loadSetting(workDir string, buntClient *buntdb.DB, logger *zap.Logger) error{
	file, err := os.Open(filepath.Join(workDir, "config", "config.yaml"))
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
	var setting models.Setting
	if err := yaml.Unmarshal(content, &setting); err != nil {
		logger.Error(fmt.Sprintf("解析默认模板文件失败: [%s]",err.Error()))
		return err
	}
	settingByte, err := json.Marshal(setting)
	if err != nil {
		logger.Error(fmt.Sprintf("序列化默认模板文件失败: [%s]",err.Error()))
		return err
	}
	if err := utils.SetValue(buntClient, "setting", string(settingByte), logger); err != nil {
		logger.Error(fmt.Sprintf("默认模板文件写入buntDB失败: [%s]",err.Error()))
		return err
	}
	return nil
}