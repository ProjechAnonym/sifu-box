package route

import (
	"io"
	"net/http"
	"path/filepath"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

// SettingUpgrade 配置升级相关路由和处理函数
// 参数:
//   group *gin.RouterGroup: Gin的路由组
//   lock *sync.Mutex: 用于并发控制的锁
func SettingUpgrade(group *gin.RouterGroup,lock *sync.Mutex) {
    // 创建升级特定资源的路由组
    route := group.Group("/upgrade")
    // 使用中间件进行令牌认证
    route.Use(middleware.TokenAuth())
    // 处理singbox升级的POST请求
    route.POST("/singbox", func(ctx *gin.Context) {
        // 获取上传的文件
        file,err := ctx.FormFile("file")
        if err != nil {
            // 如果解析文件失败，返回500错误
            ctx.JSON(http.StatusInternalServerError,gin.H{"message":"解析文件失败"})
            return
        }
        // 获取提交的地址数组
        addresses := ctx.PostFormArray("addr")
        // 打开上传的文件
        src,err := file.Open()
        if err != nil {
            // 如果打开文件失败，返回500错误
            ctx.JSON(http.StatusInternalServerError,gin.H{"message":"打开文件失败"})
            return
        }
        defer src.Close() // 确保文件在函数结束时关闭
        // 读取文件内容
        content,err := io.ReadAll(src)
        if err != nil {
            // 如果读取文件失败，记录日志并返回500错误
            utils.LoggerCaller("读取文件失败", err, 1)
            ctx.JSON(http.StatusInternalServerError,gin.H{"message":"读取文件失败"})
            return
        }
        // 执行升级工作流
        upgradeErrors := controller.UpgradeWorkflow(content,addresses,filepath.Join("/opt","singbox","sing-box"),"sing-box",lock)
        if len(upgradeErrors) == 0 {
            // 如果没有升级错误，返回200和成功消息
            ctx.JSON(http.StatusOK,gin.H{"message":true})
        }else{
            // 如果有升级错误，将错误转换为字符串数组并返回500错误
            upgradeErrorsString := make([]string,len(upgradeErrors))
            for i,upgradeErr := range upgradeErrors{
                upgradeErrorsString[i] = upgradeErr.Error()
            }
            ctx.JSON(http.StatusInternalServerError,gin.H{"message":upgradeErrorsString})
        }
    })
}