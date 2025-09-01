package main

import (
	"context"
	"fmt"
	"sifu-box/application"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/ent/template"
	"sifu-box/initial"
	"sifu-box/middleware"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
)

var config string
var dir string
var ent_client *ent.Client
var bunt_client *buntdb.DB

func init() {
	config, dir = cmd.Command()
	init_logger := initial.GetLogger(dir, "init", false)
	defer init_logger.Sync()
	ent_client = initial.InitEntdb(dir)
	bunt_client = initial.InitBuntdb()
	init_logger.Info("初始化数据库成功")
}
func main() {
	signal_chan := make(chan application.Signal, 5)
	hook_chan := make(chan application.SignalHook, 5)
	cron_chan := make(chan bool, 5)
	web_chan := make(chan bool, 5)
	exec_lock := sync.Mutex{}
	task_logger := initial.GetLogger(dir, "task", true)
	web_logger := initial.GetLogger(dir, "web", true)
	go application.ServiceControl(&signal_chan, task_logger, dir, bunt_client, &hook_chan)
	go application.HookHandle(&hook_chan, &cron_chan, &web_chan, task_logger)
	defer func() {
		web_logger.Sync()
		task_logger.Sync()
		ent_client.Close()
		bunt_client.Close()
	}()
	// if err := ent_client.Provider.Create().SetName("M78").SetRemote(true).SetPath("https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d").Exec(context.Background()); err != nil {
	// 	fmt.Println(err)
	// }
	signal_chan <- application.Signal{Operation: application.BOOT_SERVICE, Cron: true}
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
		application.Process(dir, ent_client, task_logger)
		name, err := utils.GetValue(bunt_client, initial.ACTIVE_TEMPLATE, task_logger)
		if err != nil {
			task_logger.Error(fmt.Sprintf("获取激活模板失败: [%s]", err.Error()))
			return
		} else if name == "" {
			task_logger.Error("未设置激活模板")
			return
		}
		signal_chan <- application.Signal{Cron: false, Operation: application.RELOAD_SERVICE}
		select {
		case res := <-cron_chan:
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
		task_logger.Error(fmt.Sprintf("添加定时任务失败: [%s]", err.Error()))
	}
	fmt.Println(job_id)

	application.Process(dir, ent_client, task_logger)
	name, _ := ent_client.Template.Query().Select(template.FieldName).First(context.Background())
	utils.SetValue(bunt_client, initial.ACTIVE_TEMPLATE, name.Name, task_logger)
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.Use(middleware.Logger(web_logger), middleware.Recovery(true, web_logger))
	server.Run(":8080")

	// test(ent_client)

}
func test(ent_client *ent.Client) {
	log := singbox.Log{Disabled: true}
	experiment := singbox.Experiment{Clash_api: singbox.Clash_api{External_controller: "127.0.0.1:9090", External_ui: "/ui", Secret: "123456"}}
	dns := singbox.DNS{
		Servers: []map[string]any{
			{"tag": "google", "type": "tls", "server": "8.8.8.8", "server_port": 853},
			{"tag": "cloudflare", "type": "tls", "server": "1.1.1.1", "server_port": 853},
		},
		Final:    "google",
		Strategy: "prefer_ipv4",
	}
	outbounds := []singbox.OutboundGroup{{Type: "direct", Tag: "direct"}, {Type: "selector", Tag: "selector", Providers: []string{"M78"}}, {Type: "urltest", Tag: "auto", Providers: []string{"M78"}}}
	inbounds := []map[string]any{{"tag": "tun_in", "type": "tun", "interface_name": "tun0", "mtu": 1500, "stack": "mixed", "auto_route": true, "strict_route": true, "address": []string{"172.18.0.1/30", "fdfe:dcba:9876::1/126"}}}
	route := singbox.Route{
		Default_domain_resolver: map[string]any{"server": "google"},
		Final:                   "direct",
		Rules:                   []map[string]any{{"user": []string{"bind"}, "action": "route", "outbound": "direct"}, {"port": []int{53}, "action": "hijack-dns"}, {"protocol": []string{"dns"}, "action": "hijack-dns"}, {"ip_is_private": true, "action": "route", "outbound": "direct"}, {"protocol": []string{"quic"}, "action": "reject"}},
		Rule_sets:               []singbox.Rule_set{{Type: "remote", Tag: "china-ip", Format: "binary", URL: "https://github.com/MetaCubeX/meta-rules-dat/raw/bd4354ba7f11a22883b919ac9fb9f7034fb51b31/geo/geoip/cn.srs", Download_detour: "direct", Update_interval: "1d"}},
	}
	if err := ent_client.Template.Update().Where(template.NameEQ("default")).SetDNS(dns).SetExperiment(experiment).SetInbounds(inbounds).SetRoute(route).SetOutboundGroups(outbounds).SetProviders([]string{"M78"}).SetUpdated(true).SetLog(log).Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
}
