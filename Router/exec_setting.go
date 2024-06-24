package router

import (
	"net/http"
	controller "sifu-box/Controller"
	middleware "sifu-box/Middleware"
	utils "sifu-box/Utils"
	"sync"

	"github.com/gin-gonic/gin"
)

// update_config 配置更新路由组
// 该函数通过gin框架的RouterGroup来定义一组与配置更新相关的路由处理方法
// 参数:
//   group *gin.RouterGroup: 路由组,用于定义一组具有相同前缀的路由
func update_config(group *gin.RouterGroup,lock *sync.Mutex) {
    // 定义一个子路由组,专门处理与配置更新相关的POST请求
    update_router := group.Group("/update")
    
    // 设置路由处理函数,处理POST /update/config请求
    // 该请求用于更新配置信息,通过请求体中的addr和config参数来指定更新的内容
    update_router.POST("/config", func(ctx *gin.Context) {
        // 从请求体中获取addr和config参数
        addr := ctx.PostForm("addr")
        config := ctx.PostForm("config")
        
        // 检查config参数是否为空,如果为空,则返回内部服务器错误和错误信息
        if config == "" {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "config is null"})
            return
        }
        
        // 调用controller层的Update_config方法来尝试更新配置
        // 如果更新失败,则记录错误日志并返回内部服务器错误和错误信息
        if err := controller.Update_config(addr, config,lock); err != nil {
            utils.Logger_caller("update config failed", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "update config failed"})
            return
        }
        
        // 如果更新成功,则返回成功的响应
        ctx.JSON(http.StatusOK, gin.H{"result": "success"})
    })
}

func refresh_items(group *gin.RouterGroup,lock *sync.Mutex) {
    refresh_router := group.Group("/refresh")
    refresh_router.GET("/items", func(ctx *gin.Context) {
        if err := controller.Refresh_items(lock); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "refresh items failed."})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"result": "success"})
    })
}

func check_service(group *gin.RouterGroup){
    check_router := group.Group("/check")
    check_router.POST("/service", func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        service := ctx.PostForm("service")
        status,err := controller.Check_status(url,service)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if status{
            ctx.JSON(http.StatusOK, gin.H{"message": true})
        }else{
            ctx.JSON(http.StatusOK, gin.H{"message": false})
        }
    })
}
func boot_service(group *gin.RouterGroup,lock *sync.Mutex){
    boot_router := group.Group("/boot")
    boot_router.POST("/service", func(ctx *gin.Context) {
        url := ctx.PostForm("url")
        service := ctx.PostForm("service")
        if err := controller.Boot_service(url,service,lock); err!=nil{
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}
func Setting_exec(group *gin.RouterGroup,lock *sync.Mutex){
    // 创建一个名为"setting"的子路由组,用于处理所有与设置相关的请求
    setting_router := group.Group("/execute")
    
    // 在"setting"子路由组上应用Token认证中间件,确保所有请求都需要通过认证
    setting_router.Use(middleware.Token_auth())
	update_config(setting_router,lock)
    refresh_items(setting_router,lock)
    check_service(setting_router)
    boot_service(setting_router,lock)
}