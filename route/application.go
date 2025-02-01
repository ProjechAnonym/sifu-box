package route

import (
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingHost(api *gin.RouterGroup, user *models.User, entClient *ent.Client, buntClient *buntdb.DB, singboxSetting models.Singbox, workDir string, rwLock *sync.RWMutex, execLock *sync.Mutex, scheduler *cron.Cron, jobID *cron.EntryID, logger *zap.Logger){
	host := api.Group("/application")
	host.Use(middleware.Jwt(user.PrivateKey, logger))
	host.GET("/fetch", func(ctx *gin.Context) {
		currentProvider, err := utils.GetValue(buntClient, models.CURRENTPROVIDER, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "获取当前机场失败"})
			return
		}
		currentTemplate, err := utils.GetValue(buntClient, models.CURRENTTEMPLATE, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "获取当前模板失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": map[string]string{"listen": singboxSetting.Listen, "secret": singboxSetting.Secret,"current_provider": currentProvider, "current_template": currentTemplate}})
	})
	host.POST("/set/:mode", func(ctx *gin.Context) {
		if !ctx.GetBool("admin"){
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		mode := ctx.Param("mode")
		value := ctx.PostForm("value")
		switch mode {
			case "provider":
				if err := control.SetApplication(workDir, value, mode, singboxSetting, buntClient, rwLock, execLock, logger); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			case "template":
				if err := control.SetApplication(workDir, value, mode, singboxSetting, buntClient, rwLock, execLock, logger); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			default: 
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "mode参数错误, 应为provider或template"})
				return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	host.POST("/interval", func(ctx *gin.Context) {
		if !ctx.GetBool("admin"){
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		interval := ctx.PostForm("interval")
		if err := control.SetInterval(workDir, interval, scheduler, jobID, entClient, buntClient, rwLock, execLock, singboxSetting, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}