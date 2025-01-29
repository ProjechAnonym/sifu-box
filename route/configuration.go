package route

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/models"
	"sifu-box/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func SettingConfiguration(api *gin.RouterGroup, workDir string, entClient *ent.Client, user models.User, buntClient *buntdb.DB, rwLock *sync.RWMutex, singboxSetting models.Singbox, logger *zap.Logger){
	configuration := api.Group("/configuration")
	configuration.Use(middleware.Jwt(user.PrivateKey, logger))
	configuration.GET("/fetch", func(ctx *gin.Context) {
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
			return
		}
		configuration, err :=control.Fetch(entClient, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": configuration})
	})
	configuration.DELETE("/items", func(ctx *gin.Context){
		if !ctx.GetBool("admin") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": []string{"非管理员用户"}})
			return
		}
		providers := ctx.PostFormArray("providers")
		rulesets := ctx.PostFormArray("rulesets")
		templates := ctx.PostFormArray("templates")
		errors := control.Delete(providers, rulesets, templates, workDir, buntClient, entClient, rwLock, logger)
		if errors != nil {
			errorList := make([]string, len(errors))
			for i, err := range errors {
				errorList[i] = err.Error()
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": errorList})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	})
	configuration.POST("/add", func(ctx *gin.Context) {
		conf := struct{
			Providers []models.Provider `json:"providers"`
			Rulesets []models.RuleSet `json:"rulesets"`
		}{}
		if err := ctx.ShouldBindJSON(&conf); err != nil {
			logger.Error(fmt.Sprintf("解析请求体失败: [%s]", err.Error()))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": []string{"解析请求体失败"}})
			return
		}
		errors := control.Add(conf.Providers, conf.Rulesets, entClient, buntClient, workDir, rwLock, logger)
		if errors != nil {
			errorList := make([]string, len(errors))
			for i, err := range errors {
				errorList[i] = err.Error()
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": errorList})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	configuration.POST("/files",func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			logger.Error(fmt.Sprintf("获取表单错误: [%s]", err.Error()))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": []string{"获取表单错误"}})
			return
		}
		files, ok := form.File["files"]
		if !ok {
			logger.Error("没有上传文件")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": []string{"没有上传文件"}})
		}
		var errors []error
		providers := make([]models.Provider, len(files))
		for i, file := range files {
			ext := filepath.Ext(file.Filename)
			fileHashName, err := utils.EncryptionMd5(file.Filename[:len(file.Filename) - len(ext)])
			if err != nil {
				logger.Error(fmt.Sprintf("计算'%s'哈希值失败: [%s]", file.Filename, err.Error()))
				errors = append(errors, fmt.Errorf("计算'%s'哈希值失败", file.Filename))
				continue
			}
			providers[i] = models.Provider{
				Name: file.Filename[:len(file.Filename) - len(ext)],
				Path: filepath.Join(workDir, models.STATICDIR, models.CLASHCONFIGFILE, fmt.Sprintf("%s.yaml",fileHashName)),
				Remote: false,
			}
			if err := ctx.SaveUploadedFile(file, filepath.Join(workDir, models.STATICDIR, models.CLASHCONFIGFILE, fmt.Sprintf("%s.yaml",fileHashName))); err != nil {
				logger.Error(fmt.Sprintf("保存文件失败: [%s]",err.Error()))
				errors = append(errors, fmt.Errorf("保存'%s'文件失败", file.Filename))
			}
		}
		if errors != nil {
			errorList := make([]string, len(errors))
			for i, err := range errors {
				errorList[i] = err.Error()
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": errorList})
			return
		}
		errors = append(errors, control.Add(providers, nil, entClient, buntClient, workDir, rwLock, logger)...)
		if errors != nil {
			errorList := make([]string, len(errors))
			for i, err := range errors {
				errorList[i] = err.Error()
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": errorList})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	configuration.POST("/template", func(ctx *gin.Context) {
		name := ctx.Query("name")
		if name == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": []string{"模板名称必须提供"}})
			return
		}
		template := models.Template{}
		if err := ctx.ShouldBindJSON(&template); err != nil {
			logger.Error(fmt.Sprintf("解析请求体失败: [%s]", err.Error()))
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": []string{"解析请求体失败"}})
			return
		}
		errors := control.Set(name, workDir, singboxSetting, template, buntClient, entClient, rwLock, logger)
		if errors != nil {
			errorList := make([]string, len(errors))
			for i, err := range errors {
				errorList[i] = err.Error()
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": errorList})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	configuration.GET("/recover", func(ctx *gin.Context){
		content, err := utils.GetValue(buntClient, models.DEFAULTTEMPLATEKEY, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "获取默认模板失败"})
			return
		}
		var template models.Template
		if err := yaml.Unmarshal([]byte(content), &template); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "解析默认模板失败"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message":template})
	})
}