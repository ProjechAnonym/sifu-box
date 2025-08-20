package application

import (
	"fmt"
	"strings"
	"sync"

	"github.com/codeskyblue/go-sh"
	"go.uber.org/zap"
)

func BootService(logger *zap.Logger, dir, name string, exec_lock *sync.Mutex) error {
	for {
		if exec_lock.TryLock() {
			break
		}
	}
	defer exec_lock.Unlock()
	res, err := sh.Command(fmt.Sprintf(`%s/sing-box/sing-box`, dir), "-D", fmt.Sprintf(`%s/sing-box/lib`, dir), "-c", fmt.Sprintf(`%s/config/%s.json`, dir, name), "run").CombinedOutput()
	outputs := strings.Trim(string(res), "\n")
	lines := strings.Split(outputs, "\n")
	for _, line := range lines {
		// if strings.Contains(line, "Active:") && (strings.Contains(line, "inactive") || strings.Contains(line, "activating")) {
		// 	return false, nil
		// }\
		fmt.Println(line)

	}
	if err != nil {
		logger.Error(fmt.Sprintf("执行查看状态命令失败: [%s]", err.Error()))
		return err
	}
	return nil
}
