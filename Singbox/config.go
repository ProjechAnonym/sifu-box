package singbox

import (
	"fmt"
	"io/fs"
	utils "sifu-box/Utils"

	"github.com/bitly/go-simplejson"
)
func config_merge(template string) error{
	// 获取固定信息
	log,err := utils.Get_value(template,"log")
	if err != nil{
		utils.Logger_caller("Get log failed!",err,1)
		return err
	}
	dns,err := utils.Get_value(template,"dns")
	if err != nil{
		utils.Logger_caller("Get dns failed!",err,1)
		return err
	}
	inbounds,err := utils.Get_value(template,"inbounds")
	if err != nil{
		utils.Logger_caller("Get inbounds failed!",err,1)
		return err
	}
	experimental,err := utils.Get_value(template,"experimental")
	if err != nil{
		utils.Logger_caller("Get experimental failed!",err,1)
		return err
	}
	route,err := Merge_route(template)
	if err != nil{
		utils.Logger_caller("Get route failed!",err,1)
		return err
	}
	config := simplejson.New()
	config.Set("log", log)
	config.Set("dns", dns)
	config.Set("inbounds", inbounds)
	config.Set("route", route)
	config.Set("experimental", experimental)
	config_bytes,_ := config.EncodePretty()
	err = utils.File_write(config_bytes,"E:/Myproject/sifu-box/static/1.json",[]fs.FileMode{0666,0777})
	if err != nil{
		return err
	}
	return nil
}
func Config_workflow(template string) error {
	if err := utils.Load_template(template); err != nil {
		utils.Logger_caller("load the template failed",err,1)
		return fmt.Errorf("load the %s template failed",template)
	}
	if err := utils.Load_config("Proxy"); err != nil {
		utils.Logger_caller("load the Proxy config failed",err,1)
		return fmt.Errorf("load the Proxy config failed")
	}
	config_merge(template)
	return nil
}