package middleware

import (
	"fmt"
	"net/http"
	"sifu-box/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go.uber.org/zap"
)

// Jwt是一个中间件函数, 用于验证JWT令牌并设置用户信息
// 它接受一个用于签名验证的密钥、一个BuntDB客户端和一个Zap日志记录器作为参数
// 该函数返回一个Gin的HandlerFunc, 用于在HTTP请求处理中使用
func JwtAuth(secret string, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		token, err := jwt.ParseWithClaims(authorization, &model.Jwt{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			logger.Error(fmt.Sprintf(`解析"authorization"字段失败: [%s]`, err.Error()))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": `解析"authorization"字段失败`,
			})
			return
		}
		if !token.Valid || token == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token已经失效",
			})
			return
		}
		if claims, ok := token.Claims.(*model.Jwt); !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "未知错误",
			})
			return
		} else {
			ctx.Set("admin", claims.Admin)
		}
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		admin := ctx.GetBool("admin")
		if !admin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "该账户无权限",
			})
			return
		}
	}
}
