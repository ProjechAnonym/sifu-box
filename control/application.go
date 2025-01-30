package control

import (
	"fmt"
	"sifu-box/models"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SetApplication(workDir, value, mode string, singboxSetting models.Singbox, buntClient *buntdb.DB, rwLock *sync.RWMutex, execLock *sync.Mutex, logger *zap.Logger) error{
	switch mode {
	case "provider":
		if err := utils.SetValue(buntClient, models.CURRENTPROVIDER, value, logger); err != nil{
			logger.Error(fmt.Sprintf("设置当前配置机场失败: [%s]", err.Error()))
			return err
		}
	case "template":
		if err := utils.SetValue(buntClient, models.CURRENTTEMPLATE, value, logger); err != nil{
			logger.Error(fmt.Sprintf("设置当前配置机场失败: [%s]", err.Error()))
			return err
		}
	default:
		logger.Error("模式不正确, 应设置机场或模板")
		return fmt.Errorf("模式不正确, 应设置机场或模板")
	}
	
	if err := singbox.ApplyNewConfig(workDir, singboxSetting, buntClient, rwLock, execLock, logger); err != nil{
		logger.Error(fmt.Sprintf("应用新配置失败: [%s]", err.Error()))
		return err
	}
	return nil
}
