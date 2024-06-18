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
		server_config,err := utils.Get_value("Server")
		key := server_config.(utils.Server_config).Key
		if err != nil{
			utils.Logger_caller("Get key failed!",err,1)
			ctx.Abort()
			return 
		}
		if key == header {
			ctx.Set("token",header)
			return 
		}
		ctx.JSON(http.StatusUnauthorized,gin.H{"error":"unauthorized"})
		ctx.Abort()
	}
}