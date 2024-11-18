package route

import (
	"fmt"
	"io"
	"net/http"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// SettingMigrate 配置迁移相关路由
func SettingMigrate(group *gin.RouterGroup){
    // 创建/migrate路由组
    route := group.Group("/migrate")
    // 使用Token认证中间件
    route.Use(middleware.TokenAuth())
    
    // 处理GET /migrate/export请求，用于导出信息
    route.GET("/export", func(ctx *gin.Context) {
        // 调用ExportInfo函数导出信息
        info,err := controller.ExportInfo()
        if err != nil {
            // 如果发生错误，返回500错误信息
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
        }
        // 生成当前日期格式化字符串
        currentDate := time.Now().Format("2006-01-02")
        // 构造文件名
        fileName := fmt.Sprintf("sifu-box-export-%s.yaml", currentDate)
        // 设置响应头，指示文件下载
        ctx.Header("Content-Disposition", "attachment; filename="+fileName)
        ctx.Header("Content-Type", "text/plain")
        // 返回200状态码和导出的信息
        ctx.String(http.StatusOK, info)
    })
    
    // 处理POST /migrate/import请求，用于导入信息
    route.POST("/import", func(ctx *gin.Context) {
        // 解析表单中的文件
        file, err := ctx.FormFile("file")
        if err != nil {
            // 如果解析表单失败，记录日志并返回400错误信息
            utils.LoggerCaller("解析表单失败", err, 1)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析表单失败"})
            return
        }
        
        // 打开上传的文件
        openedFile,err := file.Open()
        if err != nil {
            // 如果打开文件失败，记录日志并返回500错误信息
            utils.LoggerCaller("打开文件失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "打开文件失败"})
            return
        }
        defer openedFile.Close() // 确保文件在函数结束时关闭
        
        // 读取文件内容
        content, err := io.ReadAll(openedFile)
        if err != nil {
            // 如果读取文件内容失败，记录日志并返回500错误信息
            utils.LoggerCaller("读取文件内容失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "读取文件内容失败"})
            return
        }
        
        var importYaml models.Migrate
        // 将读取的内容解析为yaml格式
        if err := yaml.Unmarshal(content,&importYaml);err != nil{
            // 如果解析yaml失败，记录日志并返回500错误信息
            utils.LoggerCaller("转换yaml格式失败", err, 1)
            ctx.JSON(http.StatusInternalServerError, gin.H{"message": "转换yaml格式失败"})
            return
        }
        
        // 调用ImportInfo函数导入信息
        if err := controller.ImportInfo(importYaml);err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
            return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": true})
    })
}
