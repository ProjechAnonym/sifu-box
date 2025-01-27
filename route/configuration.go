package route

import (
	"fmt"
	"net/http"
	"sifu-box/control"
	"sifu-box/ent"
	"sifu-box/middleware"
	"sifu-box/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

func SettingConfiguration(api *gin.RouterGroup, workDir string, entClient *ent.Client, user models.User, buntClient *buntdb.DB, rwLock *sync.RWMutex, logger *zap.Logger){
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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "非管理员用户"})
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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "解析请求体失败"})
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
}