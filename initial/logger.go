package initial

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func GetLogger(dir, task string, debug bool) *zap.Logger {
	encoder := getEncoder()
	infoWriter := getWriter("info", task, dir)
	errorWriter := getWriter("error", task, dir)
	infoCore := zapcore.NewCore(encoder, infoWriter, zapcore.InfoLevel)
	errorCore := zapcore.NewCore(encoder, errorWriter, zapcore.ErrorLevel)
	if debug {
		consoloCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		core := zapcore.NewTee(infoCore, errorCore, consoloCore)
		logger := zap.New(core, zap.AddCaller())
		return logger
	}
	core := zapcore.NewTee(infoCore, errorCore)
	logger := zap.New(core)
	return logger
}
