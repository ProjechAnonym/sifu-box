package main

import (
	"fmt"
	"os"
	database "sifu-box/Database"
	middleware "sifu-box/Middleware"
	router "sifu-box/Router"
	singbox "sifu-box/Singbox"
	utils "sifu-box/Utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
func init(){
	if err := utils.Set_value(utils.Get_Dir(),"project-dir"); err != nil {
		fmt.Fprintln(os.Stderr,"Critical error occurred, can not set the project dir, exiting.")
		os.Exit(2)
	}
	utils.Get_core()
	utils.Load_config("Server")
	database.Get_database()
}
func main() {
	server_config,err := utils.Get_value("Server")
	if err != nil {
		fmt.Fprintln(os.Stderr,"Critical error occurred, can not get the running mode, exiting.")
		os.Exit(2)
	}
	if server_config.(utils.Server_config).Server_mode{
		gin.SetMode(gin.ReleaseMode)
		server := gin.Default()
		server.Use(middleware.Logger(),middleware.Recovery(true),cors.New(middleware.Cors()))
		api_group := server.Group("/api")
		router.Setting_server(api_group)
		router.Setting_box(api_group)
		server.Run(":8080")
	}else{
		singbox.Config_workflow([]int{1,0})
	}
	

}