package route

import (
	"fmt"
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/initial"
	"sifu-box/middleware"
	"sifu-box/model"
	"sifu-box/utils"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func SettingHosting(api *gin.RouterGroup, bunt_client *buntdb.DB, ent_client *ent.Client, work_dir string, logger *zap.Logger) {
	content, err := utils.GetValue(bunt_client, initial.USER, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("获取用户配置信息失败: [%s]", err.Error()))
		panic(fmt.Sprintf("获取用户配置信息失败: [%s]", err.Error()))
	}
	user := model.User{}
	if err := yaml.Unmarshal([]byte(content), &user); err != nil {
		logger.Error(fmt.Sprintf("序列化用户配置信息失败: [%s]", err.Error()))
		panic(fmt.Sprintf("序列化用户配置信息失败: [%s]", err.Error()))
	}
	hosting := api.Group("/files")
	hosting.GET("/list", middleware.JwtAuth(user.Key, logger), func(ctx *gin.Context) {
		files, err := control.FileList(ent_client, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "获取文件列表失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": files})
	})
	hosting.GET("/download/:path/:name", func(ctx *gin.Context) {
		path := ctx.Param("path")
		name := ctx.Param("name")
		content, err := control.FileDownload(work_dir, path, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.Header("Content-Type", "text/plain")
		ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.json"`, name))
		ctx.Data(http.StatusOK, "application/octet-stream", content)
	})
}
