package control

import (
	"errors"
	"fmt"
	"sifu-box/application"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

// OperationSingBox 对SingBox服务执行指定的操作
// operation: 要执行的操作, 支持"boot"、"reload"、"stop"、"check"
// signal_chan: 用于发送操作信号的通道
// web_chan: 用于接收操作结果的通道
// logger: 日志记录器
// 返回值1: 操作是否成功
// 返回值2: 操作结果的错误信息
func OperationSingBox(operation string, signal_chan *chan application.Signal, web_chan *chan application.ResSignal, logger *zap.Logger) (bool, error) {
	// 根据操作类型向信号通道发送对应的服务操作信号
	switch operation {
	case "boot":
		*signal_chan <- application.Signal{Cron: false, Operation: application.BOOT_SERVICE}
	case "reload":
		*signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
	case "stop":
		*signal_chan <- application.Signal{Cron: false, Operation: application.STOP_SERVICE}
	case "check":
		*signal_chan <- application.Signal{Cron: false, Operation: application.CHECK_SERVICE}
	default:
		return false, fmt.Errorf(`无效的操作"%s"`, operation)
	}

	// 等待操作结果或超时
	select {
	case res := <-*web_chan:
		if res.Content != "" {
			return res.Status, errors.New(res.Content)
		}
		return res.Status, nil
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return false, fmt.Errorf(`接收操作结果超时`)
	}
}

// RefreshFile 刷新配置文件并重新加载 sing-box 服务
// 该函数会处理工作目录下的配置更新, 并根据当前激活的模板执行相应的服务操作(如检查、重启或启动服务)
//
// 参数:
//
//	work_dir: 工作目录路径, 用于查找和处理相关配置文件
//	ent_client: 数据库客户端, 用于访问持久化数据
//	bunt_client: BuntDB 数据库客户端, 用于读取运行时键值对信息
//	signal_chan: 发送应用信号的通道, 用于通知其他协程进行特定操作
//	web_chan: 接收 Web 层响应信号的通道, 用于确认服务状态变更是否成功
//	exec_lock: 执行互斥锁, 确保同一时间只有一个实例在执行此函数逻辑
//	logger: 日志记录器, 用于输出调试及错误日志
//
// 返回值:
//
//	[]gin.H: 包含错误消息的列表, 如果一切正常则返回 nil
func RefreshFile(work_dir string, ent_client *ent.Client, bunt_client *buntdb.DB, signal_chan *chan application.Signal, web_chan *chan application.ResSignal, exec_lock *sync.Mutex, logger *zap.Logger) []gin.H {
	// 等待获取执行锁, 防止并发冲突
	for {
		if exec_lock.TryLock() {
			break
		}
	}
	defer exec_lock.Unlock()

	result := []gin.H{}

	// 处理工作目录中的配置文件, 并收集可能产生的错误
	for _, err := range application.Process(work_dir, ent_client, logger) {
		result = append(result, gin.H{"message": err.Error()})
	}

	// 获取当前激活的模板名称
	name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取激活模板失败: [%s]", err.Error()))
		result = append(result, gin.H{"message": fmt.Sprintf("获取激活模板失败: [%s]", err.Error())})
		return result
	} else if name == "" {
		return result
	}

	// 向信号通道发送检查服务状态的指令
	*signal_chan <- application.Signal{Cron: false, Operation: application.CHECK_SERVICE}

	// 等待服务状态检查的结果, 并决定是重载还是启动服务
	select {
	case res := <-*web_chan:
		if res.Status {
			*signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
		} else {
			*signal_chan <- application.Signal{Cron: false, Operation: application.BOOT_SERVICE}
		}
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		result = append(result, gin.H{"message": "查看sing-box状态超时"})
		return result
	}

	// 等待服务重载/启动操作完成后的最终反馈
	select {
	case res := <-*web_chan:
		if res.Status && len(result) == 0 {
			return nil
		} else {
			result = append(result, gin.H{"message": res.Content})
			return result
		}
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		result = append(result, gin.H{"message": "刷新配置文件重载sing-box超时"})
		return result
	}
}
