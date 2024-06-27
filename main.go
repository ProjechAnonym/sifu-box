package main

import (
	"fmt"
	"os"
	
	execute "sifu-box/Execute"
	middleware "sifu-box/Middleware"
	router "sifu-box/Router"
	singbox "sifu-box/Singbox"
	utils "sifu-box/Utils"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)
func init(){
	if err := utils.Set_value(utils.Get_Dir(),"project-dir"); err != nil {
		fmt.Fprintln(os.Stderr,"Critical error occurred, can not set the project dir, exiting.")
		os.Exit(2)
	}
	utils.Get_core()
	if err := utils.Load_config("Server");err != nil {
		utils.Logger_caller("load server config failed!",err,1)
		os.Exit(2)
	}
	if err := utils.Load_config("Proxy");err != nil{
		utils.Logger_caller("load proxy config failed!",err,1)
		os.Exit(2)
	}
    if err := utils.Load_template();err != nil{
		utils.Logger_caller("load template failed!",err,1)
		os.Exit(2)
	}
	server_config,err := utils.Get_value("Server")
	if err != nil {
		utils.Logger_caller("get server config failed!",err,1)
		os.Exit(2)
	}
	if server_config.(utils.Server_config).Server_mode{
		utils.Get_database()
	}
}
// main函数是程序的入口点
func main() {
    // 使用互斥锁来确保并发访问配置时的线程安全
    var lock sync.Mutex

    // 获取服务器配置
    server_config, err := utils.Get_value("Server")
    // 如果获取配置出错,打印错误信息并退出程序
    if err != nil {
        fmt.Fprintln(os.Stderr, "Critical error occurred, can not get the running mode, exiting.")
        os.Exit(2)
    }

    // 判断服务器是否处于服务模式
    if server_config.(utils.Server_config).Server_mode {
        // 初始化cron任务调度器
        cron_task := cron.New()
        // 每分钟执行一次配置的工作流程
        cron_id,_ := cron_task.AddFunc("@every 1m", func() {
            singbox.Config_workflow([]int{})
            var servers []utils.Server
            // 从数据库获取服务器列表
            if err := utils.Db.Find(&servers).Error; err != nil {
                utils.Logger_caller("get server list failed!", err, 1)
                return
            }
            // 获取代理配置
            proxy_config, err := utils.Get_value("Proxy")
            // 如果获取配置出错,记录错误信息
            if err != nil {
                utils.Logger_caller("get proxy config failed", err, 1)
                return
            }
            // 更新服务器组的代理配置
            execute.Group_update(servers, proxy_config.(utils.Box_config), &lock)
        })
        // 启动cron任务调度器
        cron_task.Start()

        // 设置Gin框架为发布模式
        gin.SetMode(gin.ReleaseMode)
        // 创建Gin服务器
        server := gin.Default()
        // 使用日志、恢复和跨域中间件
        server.Use(middleware.Logger(), middleware.Recovery(true), cors.New(middleware.Cors()))
        // 设置API路由
        api_group := server.Group("/api")
        router.Setting_server(api_group)
        router.Setting_box(api_group, &lock)
        router.Setting_exec(api_group, &lock,cron_task,&cron_id)
        // 启动服务器监听8080端口
        server.Run(":8080")
    } else {
        // 如果服务器不处于服务模式,只需配置工作流程
        singbox.Config_workflow([]int{})
    }
}