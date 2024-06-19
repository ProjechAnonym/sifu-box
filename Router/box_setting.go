package router

import (
	"net/http"
	controller "sifu-box/Controller"
	middleware "sifu-box/Middleware"
	utils "sifu-box/Utils"

	"github.com/gin-gonic/gin"
)
func add_items(group *gin.RouterGroup){
	add_router := group.Group("/add")
	add_router.POST("/item",func(ctx *gin.Context) {
		var config utils.Box_config
		if err := ctx.BindJSON(&config); err != nil {
			utils.Logger_caller("Marshal json failed!",err,1)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Add items failed."})
			return
		}
		if err := controller.Add_items(config);err != nil {
			utils.Logger_caller("Add items failed!",err,1)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Add items failed."})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"result": "success"})
	})

}
func Setting_box(group *gin.RouterGroup) {
	setting_router := group.Group("/config")
	setting_router.Use(middleware.Token_auth())
	add_items(setting_router)
}