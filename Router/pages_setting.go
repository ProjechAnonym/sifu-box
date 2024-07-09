package router

import (
	"net/http"
	"path/filepath"
	utils "sifu-box/Utils"

	"github.com/gin-gonic/gin"
)
func Setting_pages(server *gin.Engine) error{
	project_dir,err := utils.Get_value("project-dir")
	if err != nil{
		utils.Logger_caller("Get project dir failed",err,1)
		return err
	}
	server.StaticFS("assets",http.Dir(filepath.Join(project_dir.(string),"dist","assets")))
	server.GET("/",func(ctx *gin.Context) {
		ctx.File(filepath.Join(project_dir.(string),"dist","index.html"))
	})
	server.GET("login",func(ctx *gin.Context) {
		ctx.File(filepath.Join(project_dir.(string),"dist","index.html"))
	})
	server.GET("setting",func(ctx *gin.Context) {
		ctx.File(filepath.Join(project_dir.(string),"dist","index.html"))
	})
	return nil
}