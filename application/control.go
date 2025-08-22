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
	lines := strings.SplitSeq(outputs, "\n")
	for line := range lines {
		if strings.Contains(line, "FATAL") {
			return fmt.Errorf(`%s`, line)
		}
	}
	if err != nil {
		logger.Error(fmt.Sprintf("执行启动命令失败: [%s]", err.Error()))
		return err
	}
	return nil
}
func CheckService(logger *zap.Logger, exec_lock *sync.Mutex) (int16, error) {
	for {
		if exec_lock.TryLock() {
			break
		}
	}
	defer exec_lock.Unlock()
	res, err := sh.Command("pgrep", "-x", "sing-box").CombinedOutput()
	fmt.Println(err)
	outputs := strings.Trim(string(res), "\n")
	lines := strings.SplitSeq(outputs, "\n")
	for line := range lines {
		// if strings.Contains(line, "Active:") && (strings.Contains(line, "inactive") || strings.Contains(line, "activating")) {
		// 	return false, nil
		// }
		fmt.Println(line)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("执行查看状态命令失败: [%s]", err.Error()))
		return 0, fmt.Errorf("执行查看状态命令失败")
	}
	return 0, nil
}
