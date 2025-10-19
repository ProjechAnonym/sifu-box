package main

import (
	"fmt"
	"sifu-box/application"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/middleware"
	"sifu-box/model"
	"sifu-box/route"
	"sifu-box/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
	"gopkg.in/yaml.v3"
)

var listen string
var config_path string
var work_dir string
var ent_client *ent.Client
var bunt_client *buntdb.DB

func init() {
	config_path, work_dir, listen = cmd.Command()
	init_logger := initial.GetLogger(work_dir, "init", false)
	defer init_logger.Sync()
	ent_client = initial.InitEntdb(work_dir)
	bunt_client = initial.InitBuntdb()
	init_logger.Info("初始化数据库成功")
	initial.LoadSetting(config_path, bunt_client, init_logger)
	if err := utils.SetValue(bunt_client, initial.ACTIVE_TEMPLATE, "", init_logger); err != nil {
		init_logger.Error(fmt.Sprintf("初始化激活模板失败: [%s]", err.Error()))
		panic(fmt.Sprintf("初始化激活模板失败: [%s]", err.Error()))
	}
	if err := utils.SetValue(bunt_client, initial.OPERATION_ERRORS, "", init_logger); err != nil {
		init_logger.Error(fmt.Sprintf("初始化操作错误失败: [%s]", err.Error()))
		panic(fmt.Sprintf("初始化操作错误失败: [%s]", err.Error()))
	}
	init_logger.Info("初始化成功")
}
func main() {
	signal_chan := make(chan application.Signal, 5)
	hook_chan := make(chan application.SignalHook, 5)
	cron_chan := make(chan application.ResSignal, 5)
	web_chan := make(chan application.ResSignal, 5)
	exec_lock := sync.Mutex{}
	task_logger := initial.GetLogger(work_dir, "task", true)
	operation_logger := initial.GetLogger(work_dir, "operation", true)
	web_logger := initial.GetLogger(work_dir, "web", true)
	go application.ServiceControl(&signal_chan, task_logger, work_dir, bunt_client, &hook_chan)
	go application.HookHandle(&hook_chan, &cron_chan, &web_chan, bunt_client, task_logger)
	defer func() {
		web_logger.Sync()
		task_logger.Sync()
		operation_logger.Sync()
		ent_client.Close()
		bunt_client.Close()
	}()

	scheduler := cron.New()
	scheduler.Start()
	job_id, err := scheduler.AddFunc("* * * * *", func() {
		for {
			if exec_lock.TryLock() {
				break
			}
		}
		defer exec_lock.Unlock()
		task_logger.Info(`开始执行定时任务`)
		application.Process(work_dir, ent_client, task_logger)
		name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, task_logger)
		if err != nil {
			task_logger.Error(fmt.Sprintf("获取激活模板失败: [%s]", err.Error()))
			return
		} else if name == "" {
			task_logger.Error("未设置激活模板")
			return
		}
		signal_chan <- application.Signal{Cron: true, Operation: application.RELOAD_SERVICE}
		select {
		case res := <-cron_chan:
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
		task_logger.Error(fmt.Sprintf("添加定时任务失败: [%s]", err.Error()))
	}

	application.Process(work_dir, ent_client, task_logger)
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.Use(middleware.Logger(web_logger), middleware.Recovery(true, web_logger))
	api := server.Group("/api")

	content, err := utils.GetValue(bunt_client, initial.USER, operation_logger)
	if err != nil {
		operation_logger.Error(fmt.Sprintf("获取用户配置信息失败: [%s]", err.Error()))
		panic(fmt.Sprintf("获取用户配置信息失败: [%s]", err.Error()))
	}
	user := model.User{}
	if err := yaml.Unmarshal([]byte(content), &user); err != nil {
		operation_logger.Error(fmt.Sprintf("序列化用户配置信息失败: [%s]", err.Error()))
		panic(fmt.Sprintf("序列化用户配置信息失败: [%s]", err.Error()))
	}
	route.SettingLogin(api, &user, bunt_client, operation_logger)
	route.SettingConfiguration(api, &user, bunt_client, ent_client, work_dir, operation_logger)
	route.SettingMigrate(api, &user, ent_client, bunt_client, operation_logger)
	route.SettingExecute(api, &user, bunt_client, ent_client, work_dir, &signal_chan, &web_chan, &exec_lock, operation_logger)
	route.SettingHosting(api, &user, bunt_client, ent_client, work_dir, operation_logger)
	route.SettingApplication(api, work_dir, &user, ent_client, bunt_client, &signal_chan, &web_chan, &cron_chan, &exec_lock, scheduler, &job_id, task_logger, operation_logger)
	server.Run(listen)
}
