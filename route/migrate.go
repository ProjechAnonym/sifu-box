package route

import (
	"fmt"
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/models"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingMigrate(api *gin.RouterGroup, application *models.Application, entClient *ent.Client, buntClient *buntdb.DB, logger *zap.Logger) {
	migrate := api.Group("/migrate")
	migrate.Use(middleware.Jwt(application.Server.User.PrivateKey, logger))
	migrate.GET("export",func(ctx *gin.Context) {
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		content, err := control.Export(entClient, buntClient, application, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.Header("Content-Type", "application/octet-stream")

		// 设置Content-Disposition，提示浏览器下载文件
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "conf.yaml"))
		ctx.Data(http.StatusOK, "application/octet-stream", content)
	})
}