package route

import (
	"net/http"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/gin-gonic/gin"
)

// SettingTemplates 配置模板相关的路由
func SettingTemplates(group *gin.RouterGroup){
    // 创建/templates子路由
    route := group.Group("/templates")
    // 使用Token认证中间件
    route.Use(middleware.TokenAuth())

    // 处理模板获取请求
    route.GET("/fetch", func(ctx *gin.Context) {
        // 调用控制器方法获取模板
        templates, err := controller.GetTemplates()
        // 如果发生错误，返回500错误信息
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        // 返回200和模板信息
        ctx.JSON(http.StatusOK, gin.H{"message": templates})
    })

    // 处理模板设置请求
    route.POST("/set", func(ctx *gin.Context) {
        // 获取查询参数name
        name := ctx.Query("name")
        // 初始化模板结构体
        var template models.Template
        // 尝试从请求中解析模板信息
        if err := ctx.ShouldBindYAML(&template); err != nil {
            // 如果解析失败，记录日志并返回500错误信息
            utils.LoggerCaller("解析模板失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        // 调用控制器方法添加模板
        if err := controller.AddTemplate(name, template); err != nil {
            // 如果添加失败，返回500错误信息
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        // 返回200和成功标志
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 处理删除模板请求
    route.DELETE("delete",func(ctx *gin.Context) {
        // 获取查询参数name
        names := ctx.PostFormArray("names")
        // 调用控制器方法删除模板
        errs := controller.DeleteTemplate(names)

        if len(errs) == 0 {
            ctx.JSON(http.StatusOK, gin.H{"message": true})
        }else{
            // 获取错误信息
            errMsg := make([]string, len(errs))
            // 如果有错误，将错误转换为字符串数组并返回500错误
            if len(errs) > 0 {
                for i, err := range errs {
                    errMsg[i] = err.Error()
                }
            }
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": errMsg})
        }
        
    })

    // 处理模板刷新请求
    route.GET("/refresh", func(ctx *gin.Context){
        // 调用控制器方法刷新模板
        recoverTemplate,err := controller.RefreshTemplates()
        // 如果发生错误，返回500错误信息
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        // 返回200和恢复的模板信息
        ctx.JSON(http.StatusOK, gin.H{"message": recoverTemplate["recover"]})
    })
}