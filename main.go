package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"sifu-box/application"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/singbox"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var config string
var dir string
var ent_client *ent.Client

func init() {
	config, dir = cmd.Command()
	init_logger := initial.GetLogger(dir, "init", false)
	defer init_logger.Sync()
	ent_client = initial.InitEntdb(dir)
	init_logger.Info("初始化数据库成功")
}
func main() {
	var exec_lock = sync.Mutex{}
	taskLogger := initial.GetLogger(dir, "task", true)
	defer taskLogger.Sync()
	if err := ent_client.Provider.Create().SetName("M78").SetRemote(true).SetPath("https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d").Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
	// if err := ent_client.Provider.Create().SetName("test1").SetRemote(false).SetPath("/opt/sifubox/1.yaml").Exec(context.Background()); err != nil {
	// 	fmt.Println(err)
	// }
	if err := ent_client.Provider.Create().SetName("vless_test2").SetRemote(true).SetPath("https://raw.githubusercontent.com/Pawdroid/Free-servers/main/sub").Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
	test(ent_client)
	application.Process(dir, ent_client, taskLogger)
	// template, _ := ent_client.Template.Query().Select(template.FieldName).First(context.Background())
	fmt.Println(application.BootService(taskLogger, dir, fmt.Sprintf("%x", md5.Sum([]byte("default6"))), &exec_lock))
	application.CheckService(taskLogger, &exec_lock)
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
	if err := ent_client.Template.Create().SetName("default6").SetDNS(dns).SetExperiment(experiment).SetInbounds(inbounds).SetRoute(route).SetOutboundGroups(outbounds).SetProviders([]string{"M78"}).SetChanged(true).Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
}
