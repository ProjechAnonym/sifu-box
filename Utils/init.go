package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/huandu/go-clone"
	"github.com/spf13/viper"
)
type Box_url struct {
	Url   string `yaml:"url"`
	Proxy bool `yaml:"proxy"`
	Label string `yaml:"label"`
}
type Ruleset_value struct {
	Path   string `yaml:"path"`
	Format string `yaml:"format"`
	Type string `yaml:"type"`
	China bool `yaml:"china"`
	Update_interval string `yaml:"update_interval"`
	Download_detour string `yaml:"download_detour"`
}
type Box_ruleset struct {
	Label string `yaml:"label"`
	Value Ruleset_value `yaml:"value"`
}

type Box_config struct {
	Url     []Box_url     `yaml:"url"`
	Rule_set []Box_ruleset `yaml:"rule_set"`
}
type Cors struct {
	Origins []string `yaml:"origins"`
}
type Server_config struct {
	Cors Cors `yaml:"cors"`
	Key string `yaml:"key"`
	Server_mode bool `yaml:"server_mode"`
}
var global_vars = make(map[string]interface{})

func Get_value(keys ...string) (any, error) {
	result := clone.Clone(global_vars)
	for i, key := range keys {
		if result = result.(map[string]interface{})[key]; result == nil{
			return nil,fmt.Errorf("the key %s for level %d does not exist",key,i + 1)
		}
	}
	return result, nil
}
func Set_value(value any,keys ...string) error{
	temp_var := global_vars
	for i,key := range keys {
		if i == len(keys) - 1 {
			temp_var[key] = value
			break
		}
		if sub_map,ok := temp_var[key].(map[string]interface{}); ok{
			temp_var = sub_map
		}else{
			return fmt.Errorf("the key %s for level %d does not exist",key,i + 1)
		}	
	}
	return nil
}
func Del_key(keys ...string) error{
	temp_var := global_vars
	for i,key := range keys {
		if i == len(keys) - 1 {
			delete(temp_var,key)
			break
		}
		if sub_map,ok := temp_var[key].(map[string]interface{}); ok{
			temp_var = sub_map
		}else{
			return fmt.Errorf("the key %s for level %d does not exist",key,i + 1)
		}	
	}
	return nil
}
func Get_Dir() string {
	// base_dir := filepath.Dir(os.Args[0])
	base_dir := "E:/Myproject/sifu-box"
	return base_dir
}

func Load_config(file string) error {
	// 获取项目目录路径,获取失败直接panic退出该进程
	project_dir, err := Get_value("project-dir")
	if err != nil {
		Logger_caller(fmt.Sprintf("Get %s Dir failed!", file), err,1)
		fmt.Fprintln(os.Stderr, "Critical! Get project dictionary failed,exiting.")
		os.Exit(2)
	}
	// 读取配置文件,读取错误则panic退出该进程
	viper.SetConfigFile(filepath.Join(project_dir.(string),"config",file + ".config.yaml"))
	err = viper.ReadInConfig()
	if err != nil {
		Logger_caller(fmt.Sprintf("Read %s failed!", file), err,1)
		fmt.Fprintf(os.Stderr, "Critical! Load the %s config has failed,exiting.",file)
		os.Exit(2)
	}
	switch file {
	case "Proxy":
		var box_config Box_config
		if err := viper.Unmarshal(&box_config); err != nil {
			Logger_caller("Load proxy yaml failed",err,1)
			return err
		}
		Set_value(box_config,file)
		return nil
	case "Server":
		var server_config Server_config
		if err := viper.Unmarshal(&server_config); err != nil {
			Logger_caller("Load proxy yaml failed",err,1)
			return err
		}
		Set_value(server_config,file)
		return nil
	}
	return nil
}
func Load_template(file string) error {
	// 获取项目目录路径,获取失败直接panic退出该进程
	project_dir, err := Get_value("project-dir")
	if err != nil {
		Logger_caller(fmt.Sprintf("Get %s Dir failed!", file), err,1)
		fmt.Fprintln(os.Stderr, "Critical! Get project dictionary failed,exiting.")
		os.Exit(2)
	}
	// 读取配置文件,读取错误则panic退出该进程
	viper.SetConfigFile(filepath.Join(project_dir.(string),"template",file + ".template.yaml"))
	err = viper.ReadInConfig()
	if err != nil {
		Logger_caller(fmt.Sprintf("Read %s failed!", file), err,1)
		fmt.Fprintf(os.Stderr, "Critical! Load the %s template has failed,exiting.",file)
		os.Exit(2)
	}
	Set_value(viper.AllSettings(),file)
	return nil
}