package main

import (
	"fmt"
	"os"

	"github.com/huandu/go-clone"
	"github.com/spf13/viper"
)
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
		viper.SetConfigFile(fmt.Sprintf("%s/config/%s.config.yaml", project_dir, file))
		err = viper.ReadInConfig()
		if err != nil {
			Logger_caller(fmt.Sprintf("Read %s failed!", file), err,1)
			fmt.Fprintf(os.Stderr, "Critical! Load the %s config has failed,exiting.",file)
			os.Exit(2)
		}
		Set_value(viper.AllSettings(),file)
		return nil
}