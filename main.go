package main

import (
	"context"
	"fmt"
	"sifu-box/application"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/ent/template"
	"sifu-box/initial"
	"sifu-box/singbox"
	"sifu-box/utils"
	"sync"
	"time"

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
	signal_chan := make(chan int, 5)
	status_chan := make(chan bool, 5)
	exec_lock := sync.Mutex{}
	taskLogger := initial.GetLogger(dir, "task", true)
	defer func() {
		taskLogger.Sync()
		ent_client.Close()
	}()
	if err := ent_client.Provider.Create().SetName("M78").SetRemote(true).SetPath("https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d").Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
	scheduler := cron.New()
	scheduler.Start()
	job_id, err := scheduler.AddFunc("* * * * *", func() {
		for {
			if exec_lock.TryLock() {
				break
			}
		}
		defer exec_lock.Unlock()
		application.Process(dir, ent_client, taskLogger)
		signal_chan <- application.RELOAD_SERVICE
	})
	if err != nil {
		taskLogger.Error(fmt.Sprintf("添加定时任务失败: [%s]", err.Error()))
	}

	fmt.Println(job_id)
	if err := ent_client.Provider.Create().SetName("vless_test2").SetRemote(true).SetPath("https://raw.githubusercontent.com/Pawdroid/Free-servers/main/sub").Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
	test(ent_client)
	name, _ := ent_client.Template.Query().Select(template.FieldName).First(context.Background())
	utils.SetValue(bunt_client, initial.ACTIVE_TEMPLATE, name.Name, taskLogger)

	application.Process(dir, ent_client, taskLogger)
	go application.ServiceControl(&signal_chan, taskLogger, dir, bunt_client, &status_chan)
	signal_chan <- application.BOOT_SERVICE
	for {
		time.Sleep(time.Second * 5)
	}
}
func test(ent_client *ent.Client) {
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
	if err := ent_client.Template.Create().SetName("default").SetDNS(dns).SetExperiment(experiment).SetInbounds(inbounds).SetRoute(route).SetOutboundGroups(outbounds).SetProviders([]string{"M78"}).SetUpdated(true).Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
}
