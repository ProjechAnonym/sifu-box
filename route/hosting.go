package route

import (
	"fmt"
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingHosting(api *gin.RouterGroup, user *model.User, bunt_client *buntdb.DB, ent_client *ent.Client, work_dir string, logger *zap.Logger) {

	hosting := api.Group("/files")
	hosting.GET("/list", middleware.JwtAuth(user.Key, logger), func(ctx *gin.Context) {
		files, err := control.FileList(user.Key, ent_client, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取文件列表失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": files})
	})
	hosting.GET("/download/:name/:time/:signature/:path", func(ctx *gin.Context) {
		path := ctx.Param("path")
		name := ctx.Param("name")
		valid_time := ctx.Param("time")
		signature := ctx.Param("signature")
		content, err := control.FileDownload(work_dir, path, user.Key, name, valid_time, signature, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.Header("Content-Type", "text/plain")
		ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.json"`, name))
		ctx.Data(http.StatusOK, "application/octet-stream", content)
	})
}
