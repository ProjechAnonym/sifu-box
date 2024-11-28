package route

import (
	"net/http"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

// SettingHost 配置与主机相关的路由处理函数
// 该函数设置了三个与主机相关的端点：获取主机列表、添加主机、删除主机
// 参数 group: *gin.RouterGroup 类型,用于创建路由组
func SettingHost(group *gin.RouterGroup,lock *sync.Mutex) {
    // 创建路由组,用于管理与主机相关的请求
    route := group.Group("/host")
    // 使用中间件进行令牌认证
    route.Use(middleware.TokenAuth())

    // 设置获取主机信息的GET请求处理函数
    route.GET("fetch", func(ctx *gin.Context) {
        var hosts []models.Host
        // 从数据库中查询主机信息
        if err := utils.DiskDb.Select("url", "config", "localhost", "secret", "port","template").Find(&hosts).Error; err != nil {
            // 记录日志并返回错误信息
            utils.LoggerCaller("从数据库中获取主机失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "连接数据库失败"})
            return
        }
        // 返回成功的主机信息
        ctx.JSON(http.StatusOK, hosts)
    })

    // 设置添加主机的POST请求处理函数
    route.POST("add", func(ctx *gin.Context) {
        var content models.Host

        // 绑定JSON请求体到content变量
        if err := ctx.ShouldBindJSON(&content); err != nil {
            // 记录日志并返回绑定错误信息
            utils.LoggerCaller("反序列化json失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "反序列化失败"})
            return
        }

        // 检查主机是否为本地主机
        isLocalhost, err := controller.IsLocalhost(content.Url)
        if err != nil {
            // 返回错误信息
            ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
            return
        }
        content.Localhost = isLocalhost
        // 设置模板为默认模板
        content.Template = "default"
        // 将主机信息写入数据库
        if err := utils.DiskDb.Create(&content).Error; err != nil {
            // 记录日志并返回数据库写入错误信息
            utils.LoggerCaller("写入数据库失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "写入数据库失败"})
            return
        }

        // 返回成功添加信息
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 设置删除主机的DELETE请求处理函数
    route.DELETE("delete", func(ctx *gin.Context) {
        url := ctx.PostForm("url")

        // 从数据库中删除指定URL的主机信息
        if err := utils.DiskDb.Where("url = ?", url).Delete(&models.Host{}).Error; err != nil {
            // 记录日志并返回数据库删除错误信息
            utils.LoggerCaller("从数据库删除数据失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "无法从数据库删除数据"})
            return
        }

        // 返回成功删除信息
        ctx.JSON(http.StatusOK, gin.H{"message": true})
    })

    // 设置主机更换模板请求
    route.POST("/switch",func(ctx *gin.Context) {
        // 获取需要更换模板的主机列表以及更换的模板
        urls := ctx.PostFormArray("urls")
        template := ctx.PostForm("template")
        
        // 执行更换模板操作
        if err := utils.DiskDb.Table("hosts").Where("url IN (?)", urls).Update("template",template).Error; err != nil {
            utils.LoggerCaller("数据库查询失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "数据库查询失败"})
            return
        }

        // 更换模板需要对主机更换配置文件
        errs := controller.SwitchTemplate(template,urls,lock)
        if len(errs) == 0 {
            // 没有错误,返回true
            ctx.JSON(http.StatusOK, gin.H{"message": true})
        }else{
            // 获取错误信息
            errsMsg := make([]string,len(errs))
            for i,errMsg := range errs {
                errsMsg[i] = errMsg.Error()
            }
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": errsMsg})
        }
        
    })
}