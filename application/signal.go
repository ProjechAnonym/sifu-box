package application

import (
	"fmt"
	"sifu-box/initial"
	"sifu-box/utils"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

type SignalHook struct {
	PID  int
	Cron bool
}
type Signal struct {
	Operation int
	Cron      bool
}
type ResSignal struct {
	Status  bool
	Content string
}

// HookHandle 处理信号钩子的主函数, 监听hook_chan通道并根据处理结果向相应的响应通道发送信号
// 参数:
//
//	hook_chan: 信号钩子通道, 用于接收需要处理的信号
//	cron_chan: 定时任务响应通道, 用于发送定时任务相关的处理结果
//	web_chan: Web请求响应通道, 用于发送Web请求相关的处理结果
//	bunt_client: buntdb数据库客户端, 用于存储和读取操作错误信息
//	logger: zap日志记录器, 用于记录处理过程中的日志信息
func HookHandle(hook_chan *chan SignalHook, cron_chan *chan ResSignal, web_chan *chan ResSignal, bunt_client *buntdb.DB, logger *zap.Logger) {
	// 持续监听hook_chan通道, 处理接收到的信号
	for hook_res := range *hook_chan {
		// 从数据库中获取操作错误信息
		content, err := utils.GetValue(bunt_client, initial.OPERATION_ERRORS, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("获取操作错误信息失败: [%s]", err.Error()))
			content = fmt.Sprintf("获取操作错误信息失败: [%s]", err.Error())
		}

		// 根据信号来源类型(Cron或Web)和PID状态, 向对应的响应通道发送处理结果
		switch hook_res.Cron {
		case true:
			if hook_res.PID > 0 {
				*cron_chan <- ResSignal{Status: true, Content: content}
			} else {
				*cron_chan <- ResSignal{Status: false, Content: content}
			}
		case false:
			if hook_res.PID > 0 {
				*web_chan <- ResSignal{Status: true, Content: content}
			} else {
				*web_chan <- ResSignal{Status: false, Content: content}
			}
		}

		// 重置操作错误信息
		if err := utils.SetValue(bunt_client, initial.OPERATION_ERRORS, "", logger); err != nil {
			logger.Error(fmt.Sprintf("重置操作错误信息失败: [%s]", err.Error()))
		}
	}
}
