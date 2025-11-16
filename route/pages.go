package route

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SettingPages(server *gin.Engine, work_dir string) {
	server.StaticFS("assets", http.Dir(filepath.Join(work_dir, "static", "dist", "assets")))
	server.GET("/", func(ctx *gin.Context) {
		ctx.File(filepath.Join(work_dir, "static", "dist", "index.html"))
	})
	server.GET("/home", func(ctx *gin.Context) {
		ctx.File(filepath.Join(work_dir, "static", "dist", "index.html"))
	})
	server.GET("/setting", func(ctx *gin.Context) {
		ctx.File(filepath.Join(work_dir, "static", "dist", "index.html"))
	})
	server.GET("/login", func(ctx *gin.Context) {
		ctx.File(filepath.Join(work_dir, "static", "dist", "index.html"))
	})
}
