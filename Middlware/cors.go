package middlware

import (
	utils "sifu-box/Utils"
	"time"

	"github.com/gin-contrib/cors"
)

func Cors() cors.Config {
	origins, _ := utils.Get_value("config", "cors", "origins")
	var allow_origins = make([]string, len(origins.([]interface{})))
	for i, origin := range origins.([]interface{}) {
		allow_origins[i] = origin.(string)
	}
	cores_config := cors.Config{
		AllowOrigins:     allow_origins,
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "domain", "scheme", "Authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cores_config
}