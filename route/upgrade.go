package route

import (
	"sifu-box/middleware"

	"github.com/gin-gonic/gin"
)

func SettingUpgrade(group *gin.RouterGroup) {
	route := group.Group("/upgrade")
	route.Use(middleware.TokenAuth())
	route.POST("/singbox", func(ctx *gin.Context) {
		
	})
}