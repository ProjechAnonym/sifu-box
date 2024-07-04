package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	controller "sifu-box/Controller"
	middleware "sifu-box/Middleware"
	utils "sifu-box/Utils"

	"github.com/gin-gonic/gin"
)

func fetch_links(group *gin.RouterGroup) {
	template_router := group.Group("/fetch")
	template_router.Use(middleware.Token_auth())
	template_router.GET("/links",func(ctx *gin.Context) {
		links,err := controller.Fetch_links()
		if err != nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"fetch links failed"})
			return
		}
		ctx.JSON(http.StatusOK,links)
	})
	
}
func Setting_files(group *gin.RouterGroup) {
	setting_router := group.Group("/files")
	setting_router.GET("/:file",func(ctx *gin.Context) {
		file := ctx.Param("file")
		template := ctx.Query("template")
		token := ctx.Query("token")
		label := ctx.Query("label")
		project_dir,err := utils.Get_value("project-dir")
		if err != nil{
			utils.Logger_caller("get project dir failed",err,1)
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"get project root dir failed"})
		}
		if err := controller.Verify_link(token); err != nil{
			utils.Logger_caller("verify link failed",err,1)
			ctx.JSON(http.StatusUnauthorized,gin.H{"error":"verify link failed"})
			return
		}
		ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, label))
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.File(filepath.Join(project_dir.(string),"static",template,file))
	})
	fetch_links(setting_router)
}