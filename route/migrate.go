package route

import (
	"fmt"
	"io"
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

func SettingMigrate(api *gin.RouterGroup, ent_client *ent.Client, bunt_client *buntdb.DB, logger *zap.Logger) {
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
	migrate := api.Group("/migrate")
	migrate.Use(middleware.JwtAuth(user.Key, logger))
	migrate.GET("/export", middleware.AdminAuth(), func(c *gin.Context) {
		content, err := control.Export(ent_client, bunt_client, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.Header("Content-Type", "text/plain")
		c.Header("Content-Disposition", `attachment; filename="migrate.yaml"`)
		c.Data(http.StatusOK, "application/octet-stream", content)
	})
	migrate.POST("/import", middleware.AdminAuth(), func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("获取文件表单失败: [%s]", err.Error())})
			return
		}
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("打开文件失败: [%s]", err.Error())})
			return
		}
		defer src.Close()
		content, err := io.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("读取文件失败: [%s]", err.Error())})
			return
		}
		res, err := control.Import(content, ent_client, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusMultiStatus, res)
	})
}
