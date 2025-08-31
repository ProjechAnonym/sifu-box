package application

import (
	"crypto/md5"
	"fmt"
	"sifu-box/initial"
	"sifu-box/utils"
	"strconv"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func bootService(logger *zap.Logger, dir string, pid *int, exit *bool, bunt_client *buntdb.DB, status_chan *chan int) {
	defer func() {
		*exit = true
		*pid = 0
		*status_chan <- *pid
	}()
	name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf(`获取当前使用模板错误: [%s]`, err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf(`获取当前使用模板错误: [%s]`, err.Error()), logger)
		return
	}
	if err := sh.Command(fmt.Sprintf(`%s/sing-box/sing-box`, dir), "-D", fmt.Sprintf(`%s/sing-box/lib`, dir), "-c", fmt.Sprintf(`%s/config/%x.json`, dir, md5.Sum([]byte(name))), "run").Run(); err != nil {
		logger.Error(fmt.Sprintf("执行启动命令失败: [%s]", err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf("执行启动命令失败: [%s]", err.Error()), logger)
	}
}
func checkService(pid *int, bunt_client *buntdb.DB, logger *zap.Logger) {
	session := sh.Command("pgrep", "-x", "sing-box")
	session.SetTimeout(time.Second * 3)
	res, err := session.CombinedOutput()
	if err != nil {
		*pid = 0
		logger.Error(fmt.Sprintf(`执行查看状态命令失败: [%s]`, err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf(`执行查看状态命令失败: [%s]`, err.Error()), logger)
		return
	}
	*pid, err = strconv.Atoi(strings.Trim(string(res), "\n"))
	if err != nil {
		logger.Error(fmt.Sprintf(`转换pid数字失败: [%s]`, err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf(`转换pid数字失败: [%s]`, err.Error()), logger)
	}
}
func stopService(pid *int, bunt_client *buntdb.DB, logger *zap.Logger, status_chan *chan int) {
	if *pid <= 0 {
		*status_chan <- *pid
		logger.Error(`关闭失败, 未找到进程PID`)
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, `未找到进程, 重新检查进程状态`, logger)
		return
	}
	if err := sh.Command("kill", "-9", fmt.Sprintf("%d", *pid)).Run(); err != nil {
		*status_chan <- *pid
		logger.Error(fmt.Sprintf(`执行停止命令失败: [%s]`, err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf(`执行停止命令失败: [%s]`, err.Error()), logger)
		return
	}
}
func reloadService(pid *int, bunt_client *buntdb.DB, logger *zap.Logger) {
	if *pid <= 0 {
		*pid = 0
		logger.Error(`重载失败, 未找到进程PID`)
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, `未找到进程, 重新检查进程状态`, logger)
		return
	}
	if err := sh.Command("kill", "-HUP", fmt.Sprintf("%d", *pid)).Run(); err != nil {
		*pid = 0
		logger.Error(fmt.Sprintf(`执行重载命令失败: [%s]`, err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf(`执行重载命令失败: [%s]`, err.Error()), logger)
	}
}
func ServiceControl(operation *chan int, logger *zap.Logger, dir string, buntdb_client *buntdb.DB, status_chan *chan int) {
	singbox_pid := 0
	exit := true
	for signal := range *operation {
		switch signal {
		case BOOT_SERVICE:
			if !exit {
				continue
			}
			exit = false
			go bootService(logger, dir, &singbox_pid, &exit, buntdb_client, status_chan)
			*operation <- CHECK_SERVICE
		case CHECK_SERVICE:
			checkService(&singbox_pid, buntdb_client, logger)
			*status_chan <- singbox_pid
		case RELOAD_SERVICE:
			reloadService(&singbox_pid, buntdb_client, logger)
			*status_chan <- singbox_pid
		case STOP_SERVICE:
			stopService(&singbox_pid, buntdb_client, logger, status_chan)
		}
	}

}
