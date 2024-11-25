package route

import (
	"net/http"
	"path/filepath"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sync"

	"github.com/gin-gonic/gin"
)

func SettingUpgrade(group *gin.RouterGroup,lock *sync.Mutex) {
	route := group.Group("/upgrade")
	route.Use(middleware.TokenAuth())
	route.POST("/singbox", func(ctx *gin.Context) {
		file,err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"message":"解析文件失败"})
			return
		}
		addr := ctx.PostForm("addr")
		src,err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"message":"打开文件失败"})
			return
		}
		defer src.Close()
		if err := controller.UpgradeApp(src,filepath.Join("/opt","singbox","sing-box"),addr,"sing-box",lock); err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"message":err.Error()})
			return
		}
		ctx.JSON(http.StatusOK,gin.H{"message":true})
	})
	
}