package middleware

import (
	"net/http"
	"sifu-box/models"
	"sifu-box/utils"

	"github.com/gin-gonic/gin"
)

// TokenAuth 生成一个用于Token认证的中间件
// 该中间件通过检查请求头中的Authorization字段来验证请求是否被授权
// 如果请求未携带正确的Token,中间件将返回401状态码并终止请求处理
func TokenAuth() gin.HandlerFunc {
	
	return func(ctx *gin.Context) {
		
		// 从请求头中获取Authorization字段
		header := ctx.GetHeader("Authorization")
		
		// 如果请求头为空,返回401状态码并终止请求处理
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		// 从配置中获取服务器模式配置
		serverConfig, err := utils.GetValue("mode")
		
		// 如果获取配置失败,记录错误并终止请求处理
		if err != nil {
			utils.LoggerCaller("Get key failed!", err, 1)
			ctx.Abort()
			return
		}
		
		// 从服务器配置中提取Token
		key := serverConfig.(models.Server).Token

		// 如果请求头中的Token与配置的Token匹配,则将Token设置到上下文中并继续请求处理
		if key == header {
			
			ctx.Set("token", header)
			return
		}
		
		// 如果Token不匹配,返回401状态码并终止请求处理
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.Abort()
	}
}