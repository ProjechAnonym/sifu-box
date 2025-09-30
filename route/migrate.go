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
	migrate.GET("/fetch", func(c *gin.Context) {
		content, err := control.Migrate(ent_client, bunt_client, logger)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.Header("Content-Type", "text/plain")
		c.Header("Content-Disposition", `attachment; filename="migrate.yaml"`)
		c.Data(http.StatusOK, "application/octet-stream", content)
	})
}
