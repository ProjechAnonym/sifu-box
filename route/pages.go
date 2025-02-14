package route

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SettingPages(server *gin.Engine, workDir string) {
	server.StaticFS("assets", http.Dir(filepath.Join(workDir, "static", "dist", "assets")))
	server.GET("/",func(ctx *gin.Context) {
		ctx.File(filepath.Join(workDir, "static", "dist", "index.html"))
	})
	server.GET("/home",func(ctx *gin.Context) {
		ctx.File(filepath.Join(workDir, "static", "dist", "index.html"))
	})
	server.GET("/setting",func(ctx *gin.Context) {
		ctx.File(filepath.Join(workDir, "static", "dist", "index.html"))
	})
	server.GET("/login",func(ctx *gin.Context) {
		ctx.File(filepath.Join(workDir, "static", "dist", "index.html"))
	})
}