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

func SettingConfiguration(api *gin.RouterGroup, bunt_client *buntdb.DB, ent_client *ent.Client, work_dir string, logger *zap.Logger) {
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
		msg := control.FetchItems(ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, gin.H{"message": msg})
	})
	configuration.POST("/add/provider/:remote", middleware.AdminAuth(), func(ctx *gin.Context) {
		providers := []model.Provider{}
		res := []gin.H{}
		switch ctx.Param("remote") {
		case "remote":
			if err := ctx.BindJSON(&providers); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析JSON失败"})
				return
			}
		case "local":
			form, err := ctx.MultipartForm()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "获取文件表单失败"})
				return
			}
			files := form.File["file"]
			for _, file := range files {
				provider := model.Provider{}
				if err := provider.AutoFill(file, work_dir); err != nil {
					res = append(res, gin.H{"status": false, "message": err.Error()})
					continue
				}
				if err := ctx.SaveUploadedFile(file, provider.Path); err != nil {
					res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`保存文件"%s"失败: [%s]`, file.Filename, err.Error())})
					continue
				}
				providers = append(providers, provider)
			}
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "未标明云端或本地"})
			return
		}
		res = append(res, control.AddProvider(providers, ent_client, logger)...)
		ctx.JSON(http.StatusMultiStatus, gin.H{"message": res})
	})
	configuration.PATCH("/edit/provider", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		path := ctx.PostForm("path")
		remote := ctx.PostForm("remote") == "true"
		if err := control.EditProvider(name, path, remote, ent_client, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(`修改机场"%s"配置成功`, name)})
	})
	configuration.DELETE("/delete/provider", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostFormArray("name")
		res := control.DeleteProvider(name, ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, res)
	})
	configuration.POST("/add/ruleset/:remote", middleware.AdminAuth(), func(ctx *gin.Context) {
		rulesets := []model.Ruleset{}
		res := []gin.H{}
		switch ctx.Param("remote") {
		case "remote":
			if err := ctx.BindJSON(&rulesets); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析JSON失败"})
				return
			}
		case "local":
			form, err := ctx.MultipartForm()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "获取文件表单失败"})
				return
			}
			files := form.File["file"]
			for _, file := range files {
				ruleset := model.Ruleset{}
				if err := ruleset.AutoFill(file, work_dir); err != nil {
					res = append(res, gin.H{"status": false, "message": err.Error()})
					continue
				}
				if err := ctx.SaveUploadedFile(file, ruleset.Path); err != nil {
					res = append(res, gin.H{"status": false, "message": fmt.Sprintf(`保存文件"%s"失败: [%s]`, file.Filename, err.Error())})
					continue
				}
				rulesets = append(rulesets, ruleset)
			}
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "未标明云端或本地"})
			return
		}
		res = append(res, control.AddRuleset(rulesets, ent_client, logger)...)
		ctx.JSON(http.StatusMultiStatus, gin.H{"message": res})
	})
	configuration.PATCH("/edit/ruleset", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		path := ctx.PostForm("path")
		download_detour := ctx.PostForm("download_detour")
		update_interval := ctx.PostForm("update_interval")
		remote := ctx.PostForm("remote") == "true"
		binary := ctx.PostForm("binary") == "true"

		if err := control.EditRuleset(name, path, update_interval, download_detour, remote, binary, ent_client, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(`修改规则集"%s"成功`, name)})
	})
	configuration.DELETE("/delete/ruleset", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostFormArray("name")
		res := control.DeleteRuleset(name, ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, res)
	})
	configuration.POST("/add/template", middleware.AdminAuth(), func(ctx *gin.Context) {
		template := model.Template{}
		if err := ctx.BindJSON(&template); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "解析JSON失败"})
			return
		}
		if err := template.CheckField(); err != nil {
			logger.Error(fmt.Sprintf(`模板字段"%s"出错: [%s]`, template.Name, err.Error()))
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf(`模板字段"%s"出错: [%s]`, template.Name, err.Error())})
			return
		}
		if err := control.AddTemplate(template, ent_client, logger); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(`添加模板"%s"成功`, template.Name)})
	})
	configuration.DELETE("/delete/template", middleware.AdminAuth(), func(ctx *gin.Context) {
		name := ctx.PostFormArray("name")
		res := control.DeleteTemplate(name, work_dir, ent_client, logger)
		ctx.JSON(http.StatusMultiStatus, res)
	})
}
