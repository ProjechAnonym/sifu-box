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

func bootService(logger *zap.Logger, dir string, pid *int, name_chan *chan string, err_chan *chan error, exit_chan *chan int) {
	defer func() {
		*exit_chan <- 0
		*pid = 0
	}()
	select {
	case name := <-*name_chan:
		err := sh.Command(fmt.Sprintf(`%s/sing-box/sing-box`, dir), "-D", fmt.Sprintf(`%s/sing-box/lib`, dir), "-c", fmt.Sprintf(`%s/config/%x.json`, dir, md5.Sum([]byte(name))), "run").Run()
		if err != nil {
			logger.Error(fmt.Sprintf("执行启动命令失败: [%s]", err.Error()))
			*err_chan <- err
		}
	case <-time.After(10 * time.Second):
		logger.Error(`执行启动命令失败, 未接收到配置名称`)
		*err_chan <- errors.New("执行启动命令失败, 未接收到配置名称")
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
		*err_chan <- errors.New("未找到进程, 重新检查进程状态")
		return
	}
	sh.Command("kill", "-9", fmt.Sprintf("%d", *pid)).Run()
}

func ServiceControl(operation *chan int, logger *zap.Logger, dir string, err_chan *chan error, name_chan *chan string) {
	singbox_pid := 0
	exit_chan := make(chan int, 5)
	boot_process := false
	for {
		select {
		case op := <-*operation:
			switch op {
			case 0:
				if boot_process {
					continue
				}
				boot_process = true
				go bootService(logger, dir, &singbox_pid, name_chan, err_chan, &exit_chan)
				*operation <- 1
			case 1:
				checkService(&singbox_pid, logger, err_chan)
			}
		case exit_process := <-exit_chan:
			switch exit_process {
			case 0:
				boot_process = false
			}
		}

	}

}
