package route

import (
	"net/http"
	"path/filepath"
	"sifu-box/utils"

	"github.com/gin-gonic/gin"
)
func SettingPages(server *gin.Engine) error{
	project_dir,err := utils.GetValue("project-dir")
	if err != nil{
		utils.LoggerCaller("Get project dir failed",err,1)
		return err
	}
	server.StaticFS("assets",http.Dir(filepath.Join(project_dir.(string),"dist","assets")))
	server.GET("/",func(ctx *gin.Context) {
		ctx.File(filepath.Join(project_dir.(string),"dist","index.html"))
	})
	server.GET("login",func(ctx *gin.Context) {
		ctx.File(filepath.Join(project_dir.(string),"dist","index.html"))
	})
	server.GET("proxy",func(ctx *gin.Context) {
		ctx.File(filepath.Join(project_dir.(string),"dist","index.html"))
	})
	server.GET("setting",func(ctx *gin.Context) {
		ctx.File(filepath.Join(project_dir.(string),"dist","index.html"))
	})
	return nil
}