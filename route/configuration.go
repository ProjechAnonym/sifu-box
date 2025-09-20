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

func SettingConfiguration(api *gin.RouterGroup, bunt_client *buntdb.DB, ent_client *ent.Client, logger *zap.Logger) {
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
	configuration := api.Group("/configuration")
	configuration.Use(middleware.JwtAuth(user.Key, logger))
	configuration.GET("/fetch", func(ctx *gin.Context) {
		msg, err := control.FetchItems(ent_client, logger)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, msg)
			return
		} else {
			ctx.JSON(http.StatusOK, msg)
		}
	})
	configuration.POST("/add/provider", middleware.AdminAuth(), func(ctx *gin.Context) {
		provider := []struct {
			Name   string `json:"name"`
			Path   string `json:"path"`
			Remote bool   `json:"remote"`
		}{}
		if err := ctx.BindJSON(&provider); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析JSON失败"})
			return
		}
		res := control.AddProvider(provider, ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, gin.H{"message": res})
	})
	configuration.PATCH("/edit/provider", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		path := ctx.PostForm("path")
		remote := ctx.PostForm("remote") == "true"
		if err := control.EditProvider(name, path, remote, ent_client, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf(`修改机场"%s"失败`, name)})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(`修改机场"%s"成功`, name)})
	})
	configuration.DELETE("/delete/provider", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostFormArray("name")
		res := control.DeleteProvider(name, ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, res)
	})
	configuration.POST("/add/ruleset", middleware.AdminAuth(), func(ctx *gin.Context) {
		rulesets := []struct {
			Name           string `json:"name"`
			Path           string `json:"path"`
			Remote         bool   `json:"remote"`
			UpdateInterval string `json:"update_interval"`
			Binary         bool   `json:"binary"`
			DownloadDetour string `json:"download_detour"`
		}{}
		if err := ctx.BindJSON(&rulesets); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析JSON失败"})
			return
		}
		res := control.AddRuleset(rulesets, ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, gin.H{"message": res})
	})
	configuration.PATCH("/edit/ruleset", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		path := ctx.PostForm("path")
		remote := ctx.PostForm("remote") == "true"
		update_interval := ctx.PostForm("update_interval")
		binary := ctx.PostForm("binary") == "true"
		download_detour := ctx.PostForm("download_detour")
		if err := control.EditRuleset(name, path, update_interval, download_detour, remote, binary, ent_client, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf(`修改规则集"%s"失败`, name)})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(`修改规则集"%s"成功`, name)})
	})
	configuration.DELETE("/delete/ruleset", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostFormArray("name")
		res := control.DeleteRuleset(name, ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, res)
	})
}
