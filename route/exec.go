package route

import (
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)
func SettingExec(api *gin.RouterGroup, entClient *ent.Client, buntClient *buntdb.DB, workDir string, user *models.User, execLock *sync.Mutex, rwLock *sync.RWMutex, singboxSetting *models.Singbox, logger *zap.Logger)  {
	exec := api.Group("/exec")
	exec.Use(middleware.Jwt(user.PrivateKey, logger))
	exec.GET("/boot", func(ctx *gin.Context) {
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		if err := control.BootService(logger, singboxSetting.Commands[models.BOOTCOMMAND], execLock); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	exec.GET("/stop", func(ctx *gin.Context) {
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		if err := control.StopService(logger, singboxSetting.Commands[models.STOPCOMMAND], execLock); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	exec.GET("/restart", func(ctx *gin.Context){
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		if err := control.RestartService(logger, singboxSetting.Commands[models.RESTARTCOMMAND], execLock); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	exec.GET("/reload", func(ctx *gin.Context){
		if err := control.ReloadService(logger, singboxSetting.Commands[models.RELOADCOMMAND], execLock); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	exec.GET("/status", func(ctx *gin.Context){
		status, err := control.CheckService(false, logger, singboxSetting.Commands[models.CHECKCOMMAND], execLock)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": status})
	})
	exec.GET("/refresh", func(ctx *gin.Context){
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		errors := control.RefreshConf(entClient, buntClient, workDir, *singboxSetting, rwLock, execLock, logger)
		if errors != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": errors})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}