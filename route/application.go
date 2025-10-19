package route

import (
	"net/http"
	"sifu-box/application"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/model"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingApplication(api *gin.RouterGroup, work_dir string, user *model.User, ent_client *ent.Client, bunt_client *buntdb.DB, signal_chan *chan application.Signal, web_chan *chan application.ResSignal, cron_chan *chan application.ResSignal, exec_lock *sync.Mutex, scheduler *cron.Cron, job_id *cron.EntryID, task_logger *zap.Logger, logger *zap.Logger) {

	application := api.Group("/application")
	application.Use(middleware.JwtAuth(user.Key, logger))
	application.GET("/yacd", func(ctx *gin.Context) {
		yacd, err := control.FetchYacd(ent_client, bunt_client, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": yacd})
	})
	application.POST("/template", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		if err := control.SetTemplate(name, bunt_client, signal_chan, web_chan, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	application.POST("/interval", middleware.AdminAuth(), func(ctx *gin.Context) {
		scheduler.Remove(*job_id)
		interval := ctx.PostForm("interval")
		if interval == "" {
			ctx.JSON(http.StatusOK, gin.H{"message": "取消定时任务成功"})
			return
		}
		if err := control.SetInterval(interval, work_dir, scheduler, job_id, exec_lock, ent_client, bunt_client, signal_chan, cron_chan, task_logger, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}
