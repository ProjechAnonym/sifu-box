package router

import (
	"net/http"
	database "sifu-box/Database"
	middleware "sifu-box/Middleware"
	utils "sifu-box/Utils"

	"github.com/gin-gonic/gin"
)

func add_server(group *gin.RouterGroup) {
	add_router := group.Group("add")
	add_router.POST("/server",func(ctx *gin.Context) {
		
		var content database.Server
		// 解析前端的字符串
		if err := ctx.BindJSON(&content);err!=nil{
			utils.Logger_caller("Marshal json failed!",err,1)
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"marshal failed!"})
			return
		}
		// 补充所需数据
		if err := database.Db.Create(&content).Error; err != nil{
			utils.Logger_caller("Write msg to the database failed!",err,1)
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"write to the database failed!"})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"result":"success"})
	})
}
func Setting_server(group *gin.RouterGroup){
	setting_router := group.Group("setting")
	setting_router.Use(middleware.Token_auth())
	add_server(setting_router)
}