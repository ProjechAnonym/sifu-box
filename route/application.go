package route

import (
	"net/http"
	"sifu-box/control"
	"sifu-box/middleware"
	"sifu-box/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingHost(api *gin.RouterGroup, user *models.User, buntClient *buntdb.DB, singboxSetting models.Singbox, workDir string, rwLock *sync.RWMutex, execLock *sync.Mutex, logger *zap.Logger){
	host := api.Group("/application")
	host.Use(middleware.Jwt(user.PrivateKey, logger))
	host.POST("/:mode", func(ctx *gin.Context) {
		if !ctx.GetBool("admin"){
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		mode := ctx.Param("mode")
		value := ctx.PostForm("value")
		switch mode {
			case "provider":
				if err := control.SetApplication(workDir, value, mode, singboxSetting, buntClient, rwLock, execLock, logger); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			case "template":
				if err := control.SetApplication(workDir, value, mode, singboxSetting, buntClient, rwLock, execLock, logger); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			default: 
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "mode参数错误, 应为provider或template"})
				return
		}
		
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
}