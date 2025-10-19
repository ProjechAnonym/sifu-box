package control

import (
	"context"
	"fmt"
	"sifu-box/application"
	"sifu-box/ent"
	"sifu-box/ent/template"
	"sifu-box/initial"
	"sifu-box/model"
	"sifu-box/utils"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// FetchYacd 从数据库中获取Yacd配置信息
// 参数:
//   - ent_client: ent数据库客户端, 用于查询模板信息
//   - bunt_client: buntdb数据库客户端, 用于获取配置信息
//   - logger: zap日志记录器, 用于记录错误日志
//
// 返回值:
//   - *model.Yacd: Yacd配置信息结构体指针
//   - error: 错误信息, 如果获取或序列化失败则返回错误
func FetchYacd(ent_client *ent.Client, bunt_client *buntdb.DB, logger *zap.Logger) (*model.Yacd, error) {
	// 从buntdb中获取Yacd配置信息
	content, err := utils.GetValue(bunt_client, initial.YACD, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("获取用户配置信息失败: [%s]", err.Error())
	}

	// 将获取到的YAML格式配置信息反序列化为Yacd结构体
	yacd := model.Yacd{}
	if err := yaml.Unmarshal([]byte(content), &yacd); err != nil {
		logger.Error(fmt.Sprintf("序列化用户配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化用户配置信息失败: [%s]", err.Error())
	}

	// 获取当前激活的模板名称
	template_name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取默认模板信息失败: [%s]", err.Error()))
		return &yacd, fmt.Errorf("获取默认模板信息失败: [%s]", err.Error())
	}

	// 如果没有设置模板, 则直接返回Yacd配置
	if template_name == "" {
		return &yacd, nil
	}

	// 从ent数据库中查询模板的详细信息, 获取日志配置
	template_instance, err := ent_client.Template.Query().Where(template.NameEQ(template_name)).Select(template.FieldLog).First(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取模板信息失败: [%s]", err.Error()))
		return &yacd, fmt.Errorf("获取模板信息失败: [%s]", err.Error())
	}

	// 设置Yacd配置中的模板名称和日志开关状态
	yacd.Template = template_name
	yacd.Log = !template_instance.Log.Disabled
	return &yacd, nil
}

// SetTemplate 设置活动模板并触发相关服务操作
// name: 要设置的模板名称
// bunt_client: 数据库连接客户端, 用于存储模板配置
// signal_chan: 信号通道, 用于发送操作指令
// web_chan: 响应信号通道, 用于接收操作结果
// logger: 日志记录器
// 返回值: 操作成功返回nil, 失败返回错误信息
func SetTemplate(name string, bunt_client *buntdb.DB, signal_chan *chan application.Signal, web_chan *chan application.ResSignal, logger *zap.Logger) error {
	// 保存模板名称到数据库
	if err := utils.SetValue(bunt_client, initial.ACTIVE_TEMPLATE, name, logger); err != nil {
		return fmt.Errorf("设置模板失败: [%s]", err.Error())
	}

	// 发送检查服务信号, 并等待检查结果
	*signal_chan <- application.Signal{Cron: false, Operation: application.CHECK_SERVICE}
	select {
	case res := <-*web_chan:
		// 根据检查结果决定是重新加载服务还是启动服务
		if res.Status {
			*signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
		} else {
			*signal_chan <- application.Signal{Cron: false, Operation: application.BOOT_SERVICE}
		}
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return fmt.Errorf(`接收操作结果超时`)
	}

	// 等待服务加载完成, 并返回最终结果
	select {
	case res := <-*web_chan:
		if res.Status {
			logger.Info(fmt.Sprintf(`模板切换"%s"成功`, name))
			return nil
		}
		logger.Error(fmt.Sprintf(`载入"%s"配置失败: [%s]`, name, res.Content))
		return fmt.Errorf(`载入"%s"配置失败: [%s]`, name, res.Content)
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return fmt.Errorf(`接收操作结果超时`)
	}
}

// SetInterval 设置定时任务执行间隔及相关配置
// interval: 定时任务执行间隔表达式
// work_dir: 工作目录路径
// scheduler: cron调度器实例
// job_id: 定时任务ID指针, 用于存储创建的任务ID
// exec_lock: 执行互斥锁, 确保任务串行执行
// ent_client: 数据库客户端实例
// bunt_client: 内存数据库客户端实例
// signal_chan: 信号通道指针, 用于发送操作信号
// cron_chan: 定时任务结果通道指针, 用于接收操作结果
// task_logger: 任务专用日志记录器
// logger: 通用日志记录器
// 返回值: 错误信息, 如果添加定时任务失败则返回错误
func SetInterval(interval, work_dir string, scheduler *cron.Cron, job_id *cron.EntryID, exec_lock *sync.Mutex, ent_client *ent.Client, bunt_client *buntdb.DB, signal_chan *chan application.Signal, cron_chan *chan application.ResSignal, task_logger *zap.Logger, logger *zap.Logger) error {
	var err error
	*job_id, err = scheduler.AddFunc(interval, func() {
		// 等待获取执行锁, 确保同一时间只有一个任务在执行
		for {
			if exec_lock.TryLock() {
				break
			}
		}
		defer exec_lock.Unlock()
		logger.Info(`开始执行定时任务`)
		application.Process(work_dir, ent_client, task_logger)

		// 获取当前激活的模板配置
		name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, task_logger)
		if err != nil {
			task_logger.Error(fmt.Sprintf("获取激活模板失败: [%s]", err.Error()))
			return
		} else if name == "" {
			task_logger.Error("未设置激活模板")
			return
		}

		// 发送重载服务信号并等待执行结果
		*signal_chan <- application.Signal{Cron: true, Operation: application.RELOAD_SERVICE}
		select {
		case res := <-*cron_chan:
			if res.Status {
				task_logger.Info(`定时任务执行成功`)
			} else {
				task_logger.Error(`重载sing-box失败`)
			}
		case <-time.After(time.Second * 10):
			task_logger.Error(`接收操作结果超时`)
		}
	})
	if err != nil {
		logger.Error(fmt.Sprintf("添加定时任务失败: [%s]", err.Error()))
		return fmt.Errorf("添加定时任务失败: [%s]", err.Error())
	}
	return nil
}
