package route

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/utils"

	"github.com/gin-gonic/gin"
)

// SettingFiles 配置文件路由处理
// 该函数负责在给定的路由组中设置与文件操作相关的路由规则
// 参数 group: 用于分组设置路由的 gin 路由组实例
func SettingFiles(group *gin.RouterGroup) {
    // 创建路由规则,用于处理/files路径下的请求
    route := group.Group("/files")
    
    // 设置获取单个文件的路由规则
    route.GET("/:file", func(ctx *gin.Context) {
        // 从路径参数中获取文件名
        file := ctx.Param("file")
        // 从查询参数中获取模板名
        template := ctx.Query("template")
        // 从查询参数中获取访问令牌
        token := ctx.Query("token")
        // 从查询参数中获取文件标签
        label := ctx.Query("label")
        
        // 获取项目目录配置值
        projectDir, err := utils.GetValue("project-dir")
        if err != nil {
            // 如果获取失败,记录错误并返回内部服务器错误
            utils.LoggerCaller("获取工作目录失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取工作目录失败"})
            return
        }
        
        // 验证访问令牌
        if err := controller.VerifyLink(token); err != nil {
            // 如果令牌验证失败,返回未授权状态
            ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
            return
        }
        
        // 设置响应头,指定下载文件的名称
        ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.json"`, label))
        // 设置响应头,指定文件的MIME类型
        ctx.Header("Content-Type", "application/octet-stream")
        // 发送文件响应
        ctx.File(filepath.Join(projectDir.(string), "static", template, file))
    })
    
    // 设置获取链接列表的路由规则
    route.GET("fetch",middleware.TokenAuth(),func(ctx *gin.Context){
        // 获取链接列表
        links, err := controller.FetchLinks()
        if err != nil {
            // 如果获取失败,返回内部服务器错误
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        // 成功返回链接列表
        ctx.JSON(http.StatusOK, links)
    })
}