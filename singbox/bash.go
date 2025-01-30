package singbox

import (
	"errors"
	"fmt"
	"regexp"
	"sifu-box/models"
	"strings"
	"sync"

	"github.com/codeskyblue/go-sh"
	"go.uber.org/zap"
)

func checkService(reload bool, logger *zap.Logger, command *models.Command, execLock *sync.Mutex) (bool, error){
	if command == nil {
		logger.Error("执行命令失败, 命令不能为空")
		return false, fmt.Errorf("命令不能为空")
	}
	for {
		if execLock.TryLock(){
			break
		}
	}
	defer execLock.Unlock()
	res, err := sh.Command(command.Name, command.Args...).CombinedOutput()
	outputs := strings.Trim(string(res), "\n")
	lines := strings.Split(outputs, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Active:") && (strings.Contains(line, "inactive") || strings.Contains(line, "activating")) {
			return false, nil
		}
	}
	if err != nil {
		logger.Error(fmt.Sprintf("执行查看状态命令失败: [%s]", err.Error()))
		return false, fmt.Errorf("执行查看状态命令失败")
	}
	if reload {
		if strings.Contains(lines[len(lines) - 1], "ERROR") {
			re := regexp.MustCompile(`reload service:\s*(.*)`)
			matches := re.FindStringSubmatch(lines[len(lines) - 1])
			
			if len(matches) > 0 {return true, errors.New(matches[0])}
			return true, fmt.Errorf("重载配置遇到未知错误")
		}
	}
	return true, nil
}

func bootService(logger *zap.Logger, command *models.Command, execLock *sync.Mutex) (error) {
	if command == nil {
		logger.Error("执行命令失败, 命令不能为空")
		return fmt.Errorf("命令不能为空")
	}
	for {
		if execLock.TryLock(){
			break
		}
	}
	defer execLock.Unlock()
	result, err := sh.Command(command.Name, command.Args...).CombinedOutput()
	if err != nil {
		logger.Error(fmt.Sprintf("执行启动命令失败: [%s]", strings.Trim(string(result), "\n")))
		return fmt.Errorf("执行启动命令失败")
	}
	return nil
}

func reloadService(logger *zap.Logger, command *models.Command, execLock *sync.Mutex) (error) {
	if command == nil {
		logger.Error("执行命令失败, 命令不能为空")
		return fmt.Errorf("命令不能为空")
	}
	for {
		if execLock.TryLock(){
			break
		}
	}
	defer execLock.Unlock()
	result, err := sh.Command(command.Name, command.Args...).CombinedOutput()
	if err != nil {
		logger.Error(fmt.Sprintf("执行重载命令失败: [%s]", strings.Trim(string(result), "\n")))
		return fmt.Errorf("执行重载命令失败")
	}
	return nil
}

