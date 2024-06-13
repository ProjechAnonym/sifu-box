package main

import (
	"fmt"
	"os"
)
func init(){
	os.Setenv("GOTMPDIR", "/root/sifu-box")
	if err := Set_value(Get_Dir(),"project-dir"); err != nil {
		fmt.Fprintln(os.Stderr,"Critical error occurred, can not set the project dir, exiting.")
		os.Exit(2)
	}
	Get_core()
}
func main() {
	a,_:=Get_value("project-dir")
	Logger_caller(a.(string),nil,1)
	
}