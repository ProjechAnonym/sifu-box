package main

import (
	"fmt"
	"os"
	"sifu-box/Utils"
	"sifu-box/Database"
)
func init(){
	if err := utils.Set_value(utils.Get_Dir(),"project-dir"); err != nil {
		fmt.Fprintln(os.Stderr,"Critical error occurred, can not set the project dir, exiting.")
		os.Exit(2)
	}
	utils.Get_core()
	database.Get_database()
}
func main() {
	a,_:=utils.Get_value("project-dir")
	utils.Logger_caller(a.(string),nil,1)
	
}