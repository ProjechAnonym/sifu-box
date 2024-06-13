package middleware

import (
	"net/http"
	utils "sifu-box/Utils"

	"github.com/gin-gonic/gin"
)

func Token_auth() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == ""{
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized"})
			ctx.Abort()
			return
		}
		key,err := utils.Get_value("Server","key")
		if err != nil{
			utils.Logger_caller("Get key failed!",err,1)
			ctx.Abort()
			return 
		}
		if key.(string) == header {
			ctx.Set("token",header)
			return 
		}
		ctx.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized"})
		ctx.Abort()
	}
}