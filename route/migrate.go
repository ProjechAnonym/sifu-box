package route

import (
	"net/http"
	"sifu-box/middleware"
	"sifu-box/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SettingMigrate(api *gin.RouterGroup, user *models.User, logger *zap.Logger) {
	migrate := api.Group("/migrate")
	migrate.Use(middleware.Jwt(user.PrivateKey, logger))
	migrate.GET("export",func(ctx *gin.Context) {
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		
	})
}