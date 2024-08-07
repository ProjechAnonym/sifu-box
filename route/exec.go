package route

import (
	"net/http"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/utils"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// SettingExec 配置执行相关的路由处理函数
// - group: Gin的路由组,用于组织路由
// - lock: 同步互斥锁,用于控制并发访问
// - cronTask: 定时任务实例,用于管理定时任务
// - id: 定时任务的ID,用于标识特定的定时任务
func SettingExec(group *gin.RouterGroup, lock *sync.Mutex, cronTask *cron.Cron, id *cron.EntryID) {
    // 创建/exec路由组,用于处理执行相关的请求
    route := group.Group("/exec")
    // 使用Token认证中间件,确保只有认证过的用户才能访问接下来的路由
    route.Use(middleware.TokenAuth())

    // 处理更新配置的请求
    route.POST("/update", func(ctx *gin.Context) {
        addr := ctx.PostForm("addr")
        config := ctx.PostForm("config")
        
        // 检查配置是否为空
        if config == "" {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新配置文件为空"})
            return
        }
        
        // 尝试更新配置
        if err := controller.UpdateConfig(addr, config, lock); err != nil {
            utils.LoggerCaller("更新配置文件失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        
        // 更新成功
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 处理刷新项的请求
    route.GET("refresh", func(ctx *gin.Context) {
        if errs := controller.RefreshItems(lock); errs != nil {
            errors := make([]string, len(errs))
            for i, err := range errs {
                errors[i] = err.Error()
            }
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": errors})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 检查服务状态的请求处理
    route.POST("check", func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        service := ctx.PostForm("service")
        status, err := controller.CheckStatus(url, service)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        if status {
            ctx.JSON(http.StatusOK, gin.H{"message": true})
        } else {
            ctx.JSON(http.StatusOK, gin.H{"message": false})
        }
    })

    // 启动服务的请求处理
    route.POST("boot", func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        service := ctx.PostForm("service")
        if err := controller.BootService(url, service, lock); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 停止服务的请求处理
    route.POST("stop", func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        service := ctx.PostForm("service")
        if err := controller.StopService(url, service, lock); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 设置任务间隔的请求处理
    route.POST("interval", func(ctx *gin.Context) {
        span := ctx.PostFormArray("span")
        timeSpan := make([]int, len(span))
        var err error
        for i, num := range(span) {
            timeSpan[i], err = strconv.Atoi(num)
            if err != nil {
                ctx.JSON(http.StatusBadRequest, gin.H{
                    "message": "间隔必须是整数",
                })
                return
            }
        }
        if err := controller.SetInterval(timeSpan, cronTask, id, lock); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": false})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}