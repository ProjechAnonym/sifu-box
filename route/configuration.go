package route

import (
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SettingConfiguration(api *gin.RouterGroup, entClient *ent.Client, user models.User, logger *zap.Logger){
	configuration := api.Group("/configuration")
	configuration.Use(middleware.Jwt(user.PrivateKey, logger))
	configuration.GET("/fetch", func(ctx *gin.Context) {
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		configuration, err :=control.Fetch(entClient, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": configuration})
	})
}