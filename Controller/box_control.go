package controller

import (
	"fmt"
	"io/fs"
	"path/filepath"
	singbox "sifu-box/Singbox"
	utils "sifu-box/Utils"

	"gopkg.in/yaml.v3"
)



func Add_items(box_config utils.Box_config) error{
	project_dir,err := utils.Get_value("project-dir")
	if err != nil{
		utils.Logger_caller("Get project dir failed",err,1)
		return fmt.Errorf("get project dir failed")
	}
	if err := utils.Load_config("Proxy"); err != nil{
		utils.Logger_caller("Load proxy config failed",err,1)
		return fmt.Errorf("load proxy failed")
	}
	var index []int
	proxy_config,err := utils.Get_value("Proxy")
	proxyConfig := proxy_config.(utils.Box_config)
	var urls []utils.Box_url
	var rulesets []utils.Box_ruleset
	if err != nil{
		utils.Logger_caller("Get proxy config failed",err,1)
		return fmt.Errorf("get Proxy failed")
	}
	if len(box_config.Rule_set) == 0{
		urls = proxyConfig.Url
		urls_length := len(urls)
		if len(box_config.Url) == 0{
			return fmt.Errorf("no new links")
		}else{
			urls = append(urls,box_config.Url...)
			for i := range box_config.Url{
				index = append(index,urls_length + i)
			}
		}
	}else{
		rulesets = proxyConfig.Rule_set
		rulesets = append(rulesets, box_config.Rule_set...)
		if len(box_config.Url) != 0{
			urls = proxyConfig.Url
			urls = append(urls,box_config.Url...)
		}
	}
	var new_proxy_config utils.Box_config
	new_proxy_config.Rule_set = rulesets
	new_proxy_config.Url = urls
	new_proxy_yaml,err := yaml.Marshal(new_proxy_config)
	if err != nil{
		utils.Logger_caller("Marshal proxy config failed",err,1)
		return fmt.Errorf("marshal Proxy failed")
	}
	if err := utils.File_write(new_proxy_yaml,filepath.Join(project_dir.(string),"config","Proxy.config.yaml"),[]fs.FileMode{0644,0644});err != nil{
		utils.Logger_caller("Write Proxy config failed!",err,1)
		return err
	}
	if err := singbox.Config_workflow(index); err != nil{
		utils.Logger_caller("Config workflow failed",err,1)
		return fmt.Errorf("config workflow failed")
	}
	return nil
}