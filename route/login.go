package route

import (
	"fmt"
	"net/http"
	"sifu-box/control"
	"sifu-box/initial"
	"sifu-box/model"
	"sifu-box/utils"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func SettingLogin(api *gin.RouterGroup, bunt_client *buntdb.DB, logger *zap.Logger) {
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
	api.POST("/login/:user", func(ctx *gin.Context) {
		switch ctx.Param("user") {
		case "visitor":
			code := ctx.PostForm("code")
			if code == user.Code {
				token, err := control.Login(false, user.Key, logger)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"message": struct {
					JWT   string `json:"jwt"`
					Admin bool   `json:"admin"`
				}{JWT: token, Admin: false}})
				return
			}
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "密钥错误"})

		case "admin":
			username := ctx.PostForm("username")
			password := ctx.PostForm("password")
			if username == user.Username && password == user.Password {
				token, err := control.Login(true, user.Key, logger)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{"message": struct {
					JWT   string `json:"jwt"`
					Admin bool   `json:"admin"`
				}{JWT: token, Admin: true}})
				return
			}
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})

		default:
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "用户参数不正确"})
		}
	})
	api.GET("/verify", func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		token, admin, err := control.Verify(authorization, user.Key, logger)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": struct {
			JWT   string `json:"jwt"`
			Admin bool   `json:"admin"`
		}{JWT: token, Admin: admin}})
	})
}
