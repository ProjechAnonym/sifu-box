package control

import (
	"fmt"
	"sifu-box/application"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func OperationSingBox(operation string, signal_chan *chan application.Signal, web_chan *chan bool, logger *zap.Logger) (bool, error) {
	switch operation {
	case "boot":
		*signal_chan <- application.Signal{Cron: false, Operation: application.BOOT_SERVICE}
	case "reload":
		*signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
	case "stop":
		*signal_chan <- application.Signal{Cron: false, Operation: application.STOP_SERVICE}
	case "check":
		*signal_chan <- application.Signal{Cron: false, Operation: application.CHECK_SERVICE}
	default:
		return false, fmt.Errorf(`无效的操作"%s"`, operation)
	}

	select {
	case res := <-*web_chan:
		return res, nil
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return false, fmt.Errorf(`接收操作结果超时`)
	}
}
func RefreshFile(work_dir string, ent_client *ent.Client, bunt_client *buntdb.DB, signal_chan *chan application.Signal, web_chan *chan bool, logger *zap.Logger) []gin.H {
	res := []gin.H{}
	for _, err := range application.Process(work_dir, ent_client, logger) {
		res = append(res, gin.H{"message": err.Error()})
	}
	name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取激活模板失败: [%s]", err.Error()))
		res = append(res, gin.H{"message": fmt.Sprintf("获取激活模板失败: [%s]", err.Error())})
		return res
	} else if name == "" {
		logger.Error("未设置激活模板")
		return res
	}
	*signal_chan <- application.Signal{Cron: false, Operation: application.CHECK_SERVICE}
	select {
	case res := <-*web_chan:
		if res {
			*signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
		} else {
			*signal_chan <- application.Signal{Cron: false, Operation: application.BOOT_SERVICE}
		}
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		res = append(res, gin.H{"message": "查看sing-box状态超时"})
		return res
	}
	select {
	case res := <-*web_chan:
		if res {
			logger.Info(fmt.Sprintf(`模板切换"%s"成功`, name))
			return nil
		}
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return fmt.Errorf(`接收操作结果超时`)
	}
}
