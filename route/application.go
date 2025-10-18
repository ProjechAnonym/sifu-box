package route

import (
	"net/http"
	"sifu-box/application"
	"sifu-box/control"
	"sifu-box/middleware"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingApplication(api *gin.RouterGroup, user *model.User, bunt_client *buntdb.DB, signal_chan *chan application.Signal, web_chan *chan bool, logger *zap.Logger) {

	application := api.Group("/application")
	application.Use(middleware.JwtAuth(user.Key, logger))
	application.GET("/yacd", func(ctx *gin.Context) {
		yacd, err := control.FetchYacd(bunt_client, logger)
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
}
