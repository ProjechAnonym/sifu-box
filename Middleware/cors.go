package middleware

import (
	utils "sifu-box/Utils"
	"time"

	"github.com/gin-contrib/cors"
)

func Cors() cors.Config {
	server_config, err := utils.Get_value("Server")
	if err != nil{
		utils.Logger_caller("get server config failed",err,1)
	}
	origins := server_config.(utils.Server_config).Cors.Origins
	var allow_origins = make([]string, len(origins))
	copy(allow_origins, origins)
	cores_config := cors.Config{
		AllowOrigins:     allow_origins,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET","DELETE"},
		AllowHeaders:     []string{"Origin", "domain", "scheme", "Authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cores_config
}