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

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)
func SettingFiles(api *gin.RouterGroup, user *models.User, workDir string, entClient *ent.Client, logger *zap.Logger)  {
	api.GET("/files/fetch", middleware.Jwt(user.PrivateKey, logger),func(ctx *gin.Context){
		links, errors := control.GetFiles(user.PrivateKey, workDir, entClient, logger)
		if errors != nil {
			errorList := make([]string, len(errors))
			for i, err := range errors {
				errorList[i] = err.Error()
			}
			message := struct{
				Links map[string][]map[string]string `json:"links"`
				Errors []string `json:"errors"`}{Links: links, Errors: errorList}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": message})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": links})
	})
	api.GET("file/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		path := ctx.Query("path")
		token := ctx.Query("token")
		template := ctx.Query("template")
		verifyToken, err := utils.EncryptionMd5(user.PrivateKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "通过密钥生成token失败"})
			return
		}
		if token != verifyToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token错误"})
			return
		}
		ctx.Header("Content-Type", "application/octet-stream")

		// 设置Content-Disposition，提示浏览器下载文件
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", name))
		ctx.File(filepath.Join(workDir, models.TEMPDIR, models.SINGBOXCONFIGFILEDIR, template, path))
	})
}