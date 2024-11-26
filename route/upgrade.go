package route

import (
	"io"
	"net/http"
	"path/filepath"
	"sifu-box/controller"
	"sifu-box/middleware"
	"sifu-box/utils"
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
		addresses := ctx.PostFormArray("addr")
		src,err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,gin.H{"message":"打开文件失败"})
			return
		}
		defer src.Close()
		content,err := io.ReadAll(src)
		if err != nil {
			utils.LoggerCaller("读取文件失败", err, 1)
			ctx.JSON(http.StatusInternalServerError,gin.H{"message":"读取文件失败"})
			return
		}
		upgradeErrors := controller.UpgradeWorkflow(content,addresses,filepath.Join("/opt","singbox","sing-box"),"sing-box",lock)
		if len(upgradeErrors) == 0 {
			ctx.JSON(http.StatusOK,gin.H{"message":true})
		}else{
			upgradeErrorsString := make([]string,len(upgradeErrors))
			for i,upgradeErr := range upgradeErrors{
				upgradeErrorsString[i] = upgradeErr.Error()
			}
			ctx.JSON(http.StatusInternalServerError,gin.H{"message":upgradeErrorsString})
		}
	})
	
}