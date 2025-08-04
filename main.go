package main

import (
	"context"
	"fmt"
	"sifu-box/cmd"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/nodes"

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
	taskLogger := initial.GetLogger(dir, "task", true)
	defer taskLogger.Sync()
	if err := ent_client.Provider.Create().SetName("M78").SetRemote(true).SetPath("https://sub.m78sc.cn/api/v1/client/subscribe?token=083387dce0f02a10e8115379f9871c6d").Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
	if err := ent_client.Provider.Create().SetName("test1").SetRemote(false).SetPath("/opt/sifubox/1.yaml").Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
	if err := ent_client.Provider.Create().SetName("vless_test2").SetRemote(true).SetPath("https://raw.githubusercontent.com/Pawdroid/Free-servers/main/sub").Exec(context.Background()); err != nil {
		fmt.Println(err)
	}
	providers, _ := ent_client.Provider.Query().All(context.Background())
	nodes.Merge(providers, ent_client, taskLogger)
}
