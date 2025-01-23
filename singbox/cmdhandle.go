package singbox

import (
	"fmt"
	"sifu-box/models"
	"strings"

	"github.com/codeskyblue/go-sh"
	"go.uber.org/zap"
)

func checkService(logger *zap.Logger, command *models.Command) (bool, error){
	if command == nil {
		logger.Error("执行命令失败, 命令不能为空")
		return false, fmt.Errorf("命令不能为空")
	}
	result, err := sh.Command(command.Name, command.Args...).CombinedOutput()
	
	if err != nil {
		lines := strings.Split(string(result), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Active:") && (strings.Contains(line, "inactive") || strings.Contains(line, "activating")) {
				return false, nil
			}
		}
		logger.Error(fmt.Sprintf("执行查看状态命令失败: [%s]",err.Error()))
		return false, fmt.Errorf("执行查看状态命令失败")
	}
	return true, nil
}