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

func FetchYacd(ent_client *ent.Client, bunt_client *buntdb.DB, logger *zap.Logger) (*model.Yacd, error) {
	content, err := utils.GetValue(bunt_client, initial.YACD, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("获取用户配置信息失败: [%s]", err.Error())
	}
	yacd := model.Yacd{}
	if err := yaml.Unmarshal([]byte(content), &yacd); err != nil {
		logger.Error(fmt.Sprintf("序列化用户配置信息失败: [%s]", err.Error()))
		return nil, fmt.Errorf("序列化用户配置信息失败: [%s]", err.Error())
	}

	template_name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取默认模板信息失败: [%s]", err.Error()))
		return &yacd, fmt.Errorf("获取默认模板信息失败: [%s]", err.Error())
	}
	if template_name == "" {
		return &yacd, nil
	}
	template_instance, err := ent_client.Template.Query().Where(template.NameEQ(template_name)).Select(template.FieldLog).First(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("获取模板信息失败: [%s]", err.Error()))
		return &yacd, fmt.Errorf("获取模板信息失败: [%s]", err.Error())
	}
	yacd.Template = template_name
	yacd.Log = !template_instance.Log.Disabled
	return &yacd, nil
}
func SetTemplate(name string, bunt_client *buntdb.DB, signal_chan *chan application.Signal, web_chan *chan bool, logger *zap.Logger) error {
	if err := utils.SetValue(bunt_client, initial.ACTIVE_TEMPLATE, name, logger); err != nil {
		return fmt.Errorf("设置模板失败: [%s]", err.Error())
	}
	*signal_chan <- application.Signal{Cron: false, Operation: application.CHECK_SERVICE}
	select {
	case res := <-*web_chan:
		if res {
			*signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
		} else {
			*signal_chan <- application.Signal{Cron: false, Operation: application.BOOT_SERVICE}
		}
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return fmt.Errorf(`接收操作结果超时`)
	}
	select {
	case res := <-*web_chan:
		if res {
			logger.Info(fmt.Sprintf(`模板切换"%s"成功`, name))
			return nil
		}
		logger.Error(fmt.Sprintf(`载入"%s"配置失败, sing-box服务异常`, name))
		return fmt.Errorf(`载入"%s"配置失败, sing-box服务异常`, name)
	case <-time.After(time.Second * 10):
		logger.Error(`接收操作结果超时`)
		return fmt.Errorf(`接收操作结果超时`)
	}
}
func SetInterval(interval, work_dir string, scheduler *cron.Cron, job_id *cron.EntryID, exec_lock *sync.Mutex, ent_client *ent.Client, bunt_client *buntdb.DB, signal_chan *chan application.Signal, cron_chan *chan bool, task_logger *zap.Logger, logger *zap.Logger) error {
	var err error
	*job_id, err = scheduler.AddFunc(interval, func() {
		for {
			if exec_lock.TryLock() {
				break
			}
		}
		defer exec_lock.Unlock()
		logger.Info(`开始执行定时任务`)
		application.Process(work_dir, ent_client, task_logger)
		name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, task_logger)
		if err != nil {
			task_logger.Error(fmt.Sprintf("获取激活模板失败: [%s]", err.Error()))
			return
		} else if name == "" {
			task_logger.Error("未设置激活模板")
			return
		}
		*signal_chan <- application.Signal{Cron: true, Operation: application.RELOAD_SERVICE}
		select {
		case res := <-*cron_chan:
			if res {
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
