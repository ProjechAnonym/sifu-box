package utils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func get_encoder() zapcore.Encoder {
	encoder_config := zap.NewProductionEncoderConfig()
	encoder_config.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder_config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoder_config)
}

func get_writer(level string) zapcore.WriteSyncer{
	project_dir,_ := Get_value("project-dir")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/log/sifu-box-%s.log",project_dir.(string),level),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     15,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Get_core(){
	encoder := get_encoder()
	info_writer := get_writer("info")
	error_writer := get_writer("error")
	info_core := zapcore.NewCore(encoder,info_writer,zapcore.InfoLevel)
	error_core := zapcore.NewCore(encoder,error_writer,zapcore.ErrorLevel)
	consolo_core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout),zapcore.DebugLevel)
	core := zapcore.NewTee(info_core,error_core,consolo_core)
	my_logger := zap.New(core)
	zap.ReplaceGlobals(my_logger)
}
func Logger_caller(msg string,err error,skip int){
	if err != nil{
		_,file,line,_ := runtime.Caller(skip)
		zap.L().Error(msg,zap.String("caller",fmt.Sprintf("%s:%d",file,line)),zap.Error(err))
	}else{
		zap.L().Info(msg)
	}
}
