package middleware

import (
	"net/http"
	utils "sifu-box/Utils"

	"github.com/gin-gonic/gin"
)

// Token_auth 返回一个gin.HandlerFunc,用于验证请求中的令牌
// 该中间件检查请求头中的Authorization字段,如果令牌无效或缺失,则返回未授权的响应
func Token_auth() gin.HandlerFunc {
	// 返回一个闭包函数,作为gin.HandlerFunc
	return func(ctx *gin.Context) {
		// 尝试获取Authorization头字段
		header := ctx.GetHeader("Authorization")
		// 如果头字段为空,表示令牌不存在,返回未授权响应并终止当前请求
		if header == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		// 尝试获取服务器配置,特别是其中的键值,用于令牌验证
		server_config, err := utils.Get_value("Server")
		// 如果获取配置失败,记录错误并终止当前请求
		if err != nil {
			utils.Logger_caller("Get key failed!", err, 1)
			ctx.Abort()
			return
		}
		// 从配置中提取键值
		key := server_config.(utils.Server_config).Key

		// 比较请求头中的令牌和配置中的键值
		if key == header {
			// 如果令牌有效,设置token到上下文中,供后续使用
			ctx.Set("token", header)
			return
		}
		// 如果令牌无效,返回未授权响应并终止当前请求
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		ctx.Abort()
	}
}