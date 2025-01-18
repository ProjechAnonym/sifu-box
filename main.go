package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sifu-box/models"

	"github.com/natefinch/lumberjack"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

var workflowLogger *zap.Logger
var buntClient *buntdb.DB
var workDir string
func init() {
	var err error
	workDir = getWorkDir()
	workflowLogger = getLogger(workDir, "workflow")
	buntClient, err = buntdb.Open(":memory:")
	if err != nil {
		workflowLogger.Error(fmt.Sprintf("连接Buntdb数据库失败: [%s]",err.Error()))
		panic(err)
	}
	
	
	file,_ := os.Open(filepath.Join(workDir, "static", "default.template.yaml"))
	defer file.Close()
	content, _ := io.ReadAll(file)
	var template models.Template
	yaml.Unmarshal(content, &template)
	
	
	a,_ := json.MarshalIndent(template, "", "  ")
	fmt.Println(string(a))
}

func main() {
	defer func() {
		workflowLogger.Sync()
		buntClient.Close()
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
		Filename:   fmt.Sprintf("%s/logs/sifu-stock-%s-%s.log", workDir, task, level),
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

