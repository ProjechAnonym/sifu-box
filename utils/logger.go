package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriter(level string) zapcore.WriteSyncer {
	projectDir, _ := GetValue("project-dir")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/log/sifu-box-%s.log", projectDir.(string), level),
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GetCore() {
	encoder := getEncoder()
	infoWriter := getWriter("info")
	errorWriter := getWriter("error")
	infoCore := zapcore.NewCore(encoder, infoWriter, zapcore.InfoLevel)
	errorCore := zapcore.NewCore(encoder, errorWriter, zapcore.ErrorLevel)
	consoloCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout),zapcore.DebugLevel)
	core := zapcore.NewTee(infoCore, errorCore,consoloCore)
	myLogger := zap.New(core)
	zap.ReplaceGlobals(myLogger)
}
func LoggerCaller(msg string, err error, skip int) {
	if err != nil {
		_, file, line, _ := runtime.Caller(skip)
		zap.L().Error(msg, zap.String("caller", fmt.Sprintf("%s:%d", file, line)), zap.Error(err))
	} else {
		zap.L().Info(msg)
	}
}
