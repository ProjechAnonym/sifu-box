package main

import (
	"sifu-box/cmd"
	"sifu-box/initial"

	_ "github.com/mattn/go-sqlite3"
)

var config string
var dir string

func init() {
	config, dir = cmd.Command()
	init_logger := initial.GetLogger(dir, "init", false)
	defer init_logger.Sync()
	initial.InitEntdb(dir)
	init_logger.Info("初始化数据库成功")
}
func main() {
	taskLogger := initial.GetLogger(dir, "task", true)
	defer taskLogger.Sync()

}
