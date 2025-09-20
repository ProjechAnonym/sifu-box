package application

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"sifu-box/initial"
	"sifu-box/utils"
	"strconv"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func boot(cron bool, work_dir string, pid *int, exit *bool, bunt_client *buntdb.DB, hook_chan *chan SignalHook, logger *zap.Logger) {
	defer func() {
		*exit = true
		*pid = 0
		*hook_chan <- SignalHook{Cron: cron, PID: *pid}
	}()
	name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf(`获取当前使用模板错误: [%s]`, err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf(`获取当前使用模板错误: [%s]`, err.Error()), logger)
		return
	}
	singbox_path := filepath.Join(work_dir, "sing-box", "sing-box")
	singbox_lib := filepath.Join(work_dir, "sing-box", "lib")
	singbox_config := filepath.Join(work_dir, "sing-box", "config", fmt.Sprintf("%x.json", md5.Sum([]byte(name))))
	if err := sh.Command(singbox_path, "-D", singbox_lib, "-c", singbox_config, "run").Run(); err != nil {
		logger.Error(fmt.Sprintf("执行启动命令失败: [%s]", err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf("执行启动命令失败: [%s]", err.Error()), logger)
	}
}
func check(pid *int, bunt_client *buntdb.DB, logger *zap.Logger) {
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
func stop(pid *int, cron bool, bunt_client *buntdb.DB, logger *zap.Logger, hook_chan *chan SignalHook) {
	if *pid <= 0 {
		*hook_chan <- SignalHook{Cron: cron, PID: *pid}
		logger.Error(`关闭失败, 未找到进程PID`)
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, `未找到进程, 重新检查进程状态`, logger)
		return
	}
	if err := sh.Command("kill", "-9", fmt.Sprintf("%d", *pid)).Run(); err != nil {
		*hook_chan <- SignalHook{Cron: cron, PID: *pid}
		logger.Error(fmt.Sprintf(`执行停止命令失败: [%s]`, err.Error()))
		utils.SetValue(bunt_client, initial.OPERATION_ERRORS, fmt.Sprintf(`执行停止命令失败: [%s]`, err.Error()), logger)
		return
	}
}
func reload(pid *int, bunt_client *buntdb.DB, logger *zap.Logger) {
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
func ServiceControl(operation *chan Signal, logger *zap.Logger, work_dir string, buntdb_client *buntdb.DB, hook_chan *chan SignalHook) {
	singbox_pid := 0
	exit := true
	for signal := range *operation {
		switch signal.Operation {
		case BOOT_SERVICE:
			if !exit {
				continue
			}
			exit = false
			go boot(signal.Cron, work_dir, &singbox_pid, &exit, buntdb_client, hook_chan, logger)
			*operation <- Signal{Cron: signal.Cron, Operation: CHECK_SERVICE}
		case CHECK_SERVICE:
			check(&singbox_pid, buntdb_client, logger)
			*hook_chan <- SignalHook{Cron: signal.Cron, PID: singbox_pid}
		case RELOAD_SERVICE:
			reload(&singbox_pid, buntdb_client, logger)
			*hook_chan <- SignalHook{Cron: signal.Cron, PID: singbox_pid}
		case STOP_SERVICE:
			stop(&singbox_pid, signal.Cron, buntdb_client, logger, hook_chan)
		}
	}
}
