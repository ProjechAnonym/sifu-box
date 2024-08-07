package route

import (
	"net/http"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/gin-gonic/gin"
)

// SettingHost 配置与主机相关的路由处理函数
// 该函数设置了三个与主机相关的端点：获取主机列表、添加主机、删除主机
// 参数 group: *gin.RouterGroup 类型,用于创建路由组
func SettingHost(group *gin.RouterGroup) {
    // 创建路由组,用于管理与主机相关的请求
    route := group.Group("/host")
    // 使用中间件进行令牌认证
    route.Use(middleware.TokenAuth())

    // 设置获取主机信息的GET请求处理函数
    route.GET("fetch", func(ctx *gin.Context) {
        var hosts []models.Host
        // 从数据库中查询主机信息
        if err := utils.DiskDb.Select("url", "config", "localhost", "secret", "port").Find(&hosts).Error; err != nil {
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
}