package route

import (
	"net/http"
	"sifu-box/controller"
	"sifu-box/middleware"

	"github.com/gin-gonic/gin"
)

func SettingTemplates(group *gin.RouterGroup){
	route := group.Group("/templates")
	route.Use(middleware.TokenAuth())
	route.GET("/fetch", func(ctx *gin.Context) {
		templates,err := controller.GetTemplates()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": templates})
	})
}