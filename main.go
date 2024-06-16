package main

import (
	"fmt"
	"os"
	database "sifu-box/Database"
	singbox "sifu-box/Singbox"
	utils "sifu-box/Utils"
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
	singbox.Config_workflow("Gateway",false)
	// gin.SetMode(gin.ReleaseMode)
	// server := gin.Default()
	// server.Use(middleware.Logger(),middleware.Recovery(true),cors.New(middleware.Cors()))
	// api_group := server.Group("/api")
	// router.Setting_server(api_group)
	// server.Run()
}