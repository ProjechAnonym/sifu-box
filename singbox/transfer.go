package singbox

import (
	"fmt"
	"os"
	"path/filepath"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"

	"go.uber.org/zap"
)

func backupConfig(singboxConfigPath, workDir string, logger *zap.Logger) error {
	originalConfig, err := utils.ReadFile(singboxConfigPath)
	if err != nil {
		logger.Error(fmt.Sprintf("读取配置文件失败: [%s]", err.Error()))
		return fmt.Errorf("读取配置文件失败")
	}
	if err := utils.WriteFile(filepath.Join(workDir, models.TEMPDIR, models.BACKUPDIR, models.SINGBOXBACKUPCONFIGFILE), originalConfig, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644); err != nil {
		logger.Error(fmt.Sprintf("备份配置文件失败: [%s]", err.Error()))
		return fmt.Errorf("备份配置文件失败")
	}
	return nil
}

func recoverConfig(singboxConfigPath, workDir string, logger *zap.Logger) error {
	backupConfig, err := utils.ReadFile(filepath.Join(workDir, models.TEMPDIR, models.BACKUPDIR, models.SINGBOXBACKUPCONFIGFILE))
	if err != nil {
		logger.Error(fmt.Sprintf("读取备份配置文件失败: [%s]", err.Error()))
		return fmt.Errorf("读取备份配置文件失败")
	}
	if err := utils.WriteFile(singboxConfigPath, backupConfig, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644); err != nil {
		logger.Error(fmt.Sprintf("恢复配置文件失败: [%s]", err.Error()))
		return fmt.Errorf("恢复配置文件失败")
	}
	return nil
}

func transferConfig(singboxConfigPath, newConfigPath string, rwLock *sync.RWMutex, logger *zap.Logger) error {
	rwLock.RLock()
	defer rwLock.RUnlock()
	newConfig, err := utils.ReadFile(newConfigPath)
	if err != nil {
		logger.Error(fmt.Sprintf("读取新配置文件失败: [%s]", err.Error()))
		return fmt.Errorf("读取新配置文件失败")
	}
	if err := utils.WriteFile(singboxConfigPath, newConfig, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644); err != nil {
		logger.Error(fmt.Sprintf("写入配置文件失败: [%s]", err.Error()))
		return fmt.Errorf("写入配置文件失败")
	}
	return nil
}