package route

import (
	"fmt"
	"io"
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingMigrate(api *gin.RouterGroup, privateKey, workDir string, singboxSetting models.Singbox, rwLock *sync.RWMutex, execLock *sync.Mutex, entClient *ent.Client, buntClient *buntdb.DB, scheduler *cron.Cron, jobID *cron.EntryID, logger *zap.Logger) {
	migrate := api.Group("/migrate")
	migrate.Use(middleware.Jwt(privateKey, logger))
	migrate.GET("export",func(ctx *gin.Context) {
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		content, err := control.Export(entClient, buntClient, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// 设置Content-Disposition，提示浏览器下载文件
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "conf.yaml"))
		ctx.Data(http.StatusOK, "application/yaml", content)
	})
	migrate.POST("/import", func(ctx *gin.Context){
		file, err := ctx.FormFile("file")
		if err != nil {
			logger.Error(fmt.Sprintf("获取表单文件错误: [%s]", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取表单文件错误"})
			return 
		}
		content, err := file.Open()
		if err != nil {
			logger.Error(fmt.Sprintf("打开表单文件错误: [%s]", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "打开表单文件错误"})
			return 
		}
		defer content.Close()
		contentByte, err := io.ReadAll(content)
		if err != nil {
			logger.Error(fmt.Sprintf("读取表单文件错误: [%s]", err.Error()))
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "读取表单文件错误"})
			return
		}
		if err := control.Import(contentByte, workDir, singboxSetting, entClient, buntClient, scheduler, jobID, execLock, rwLock, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}