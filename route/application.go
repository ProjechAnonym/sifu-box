package route

import (
	"context"
	"fmt"
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/ent/template"
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
		information := struct{
			Listen string `json:"listen"`
			Secret string `json:"secret"`
			CurrentProvider string `json:"current_provider"`
			CurrentTemplate string `json:"current_template"`
			Log bool `json:"log"`
			Error string `json:"error"`
		}{
			Listen: singboxSetting.Listen,
			Secret: singboxSetting.Secret,
		}
		
		currentProvider, err := utils.GetValue(buntClient, models.CURRENTPROVIDER, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("获取当前配置机场失败: [%s]", err.Error()))
			information.Error = "获取当前配置机场失败"
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": information})
			return
		}
		information.CurrentProvider = currentProvider
		currentTemplate, err := utils.GetValue(buntClient, models.CURRENTTEMPLATE, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("获取当前配置模板失败: [%s]", err.Error()))
			information.Error = "获取当前配置模板失败"
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": information})
			return
		}
		information.CurrentTemplate = currentTemplate
		template, err := entClient.Template.Query().Select(template.FieldContent).Where(template.NameEQ(currentTemplate)).First(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("获取模板失败: [%s]", err.Error()))
			information.Error = "获取模板失败"
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": information})
			return
		}
		information.Log = !template.Content.Log.Disabled
		ctx.JSON(http.StatusOK, gin.H{"message": information})
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