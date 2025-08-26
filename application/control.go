package application

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"go.uber.org/zap"
)

func bootService(logger *zap.Logger, dir string, pid *int, name string, err_chan *chan error, exit *bool) {
	defer func() {
		*exit = true
		*pid = 0
	}()
	err := sh.Command(fmt.Sprintf(`%s/sing-box/sing-box`, dir), "-D", fmt.Sprintf(`%s/sing-box/lib`, dir), "-c", fmt.Sprintf(`%s/config/%x.json`, dir, md5.Sum([]byte(name))), "run").Run()
	if err != nil {
		logger.Error(fmt.Sprintf("执行启动命令失败: [%s]", err.Error()))
		*err_chan <- err
	}
}
func checkService(pid *int, logger *zap.Logger, err_chan *chan error) {
	session := sh.Command("pgrep", "-x", "sing-box")
	session.SetTimeout(time.Second * 3)
	res, err := session.CombinedOutput()
	if err != nil {
		*pid = 0
		logger.Error(fmt.Sprintf("执行查看状态命令失败: [%s]", err.Error()))
		*err_chan <- err
	}
	*pid, err = strconv.Atoi(strings.Trim(string(res), "\n"))
	if err != nil {
		*pid = 0
		logger.Error(fmt.Sprintf("转换pid数字失败: [%s]", err.Error()))
		*err_chan <- err
	}
}
func stopService(pid *int, logger *zap.Logger, err_chan *chan error) {
	if *pid <= 0 {
		logger.Error("关闭失败, 未找到进程PID")
		*err_chan <- errors.New("未找到进程, 重新检查进程状态")
		return
	}
	if err := sh.Command("kill", "-9", fmt.Sprintf("%d", *pid)).Run(); err != nil {
		logger.Error(fmt.Sprintf("执行停止命令失败: [%s]", err.Error()))
		*err_chan <- err
	}
}
func reloadService(pid *int, logger *zap.Logger, err_chan *chan error) {
	if *pid <= 0 {
		logger.Error("重载失败, 未找到进程PID")
		*err_chan <- errors.New("未找到进程, 重新检查进程状态")
		return
	}
	if err := sh.Command("kill", "-HUP", fmt.Sprintf("%d", *pid)).Run(); err != nil {
		logger.Error(fmt.Sprintf("执行停止命令失败: [%s]", err.Error()))
		*err_chan <- err
	}
}
func ServiceControl(operation *chan int, logger *zap.Logger, dir string, err_chan *chan error, name_chan *chan string, status_chan *chan bool) {
	singbox_pid := 0
	exit := true
	for {
		select {
		case op := <-*operation:
			switch op {
			case CHECK_SERVICE:
				checkService(&singbox_pid, logger, err_chan)
				fmt.Println(singbox_pid)
			case RELOAD_SERVICE:
				reloadService(&singbox_pid, logger, err_chan)
			case STOP_SERVICE:
				stopService(&singbox_pid, logger, err_chan)
			}
		case name := <-*name_chan:
			if !exit {
				continue
			}
			exit = false
			go bootService(logger, dir, &singbox_pid, name, err_chan, &exit)
			*operation <- CHECK_SERVICE
		}

	}
}
