package route

import (
	"net/http"
	"sifu-box/application"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingExecute(api *gin.RouterGroup, user *model.User, bunt_client *buntdb.DB, ent_client *ent.Client, work_dir string, signal_chan *chan application.Signal, web_chan *chan bool, logger *zap.Logger) {
	execute := api.Group("/execute")
	execute.Use(middleware.JwtAuth(user.Key, logger))
	execute.GET("/:operation", middleware.AdminAuth(), func(ctx *gin.Context) {
		operation := ctx.Param("operation")
		res, err := control.OperationSingBox(operation, signal_chan, web_chan, logger)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": res})
	})
	execute.GET("/refresh", func(ctx *gin.Context) {})
}
