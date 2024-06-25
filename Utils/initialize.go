package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/huandu/go-clone"
	"github.com/spf13/viper"
)

var global_vars = make(map[string]interface{})
var mu sync.RWMutex
func Get_value(keys ...string) (any, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := clone.Clone(global_vars)
	for i, key := range keys {
		if result = result.(map[string]interface{})[key]; result == nil{
			return nil,fmt.Errorf("the key %s for level %d does not exist",key,i + 1)
		}
	}
	return result, nil
}
func Set_value(value any,keys ...string) error{
	mu.Lock()
	defer mu.Unlock()
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
	mu.Lock()
	defer mu.Unlock()
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
	base_dir := "/root/sifu-box"
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
func Load_template() error {
	// 获取项目目录路径,获取失败直接panic退出该进程
	project_dir, err := Get_value("project-dir")
	if err != nil {
		Logger_caller("get project dir failed", err,1)
		return err
	}
	// 打开目录
	template_dir, err := os.Open(filepath.Join(project_dir.(string),"template"))
	if err != nil {
		Logger_caller("failed to open template directory", err,1)
		return err
	}
	defer template_dir.Close()

	// 读取目录条目
	entries, err := template_dir.ReadDir(-1) // -1 表示读取所有条目
	if err != nil {
		Logger_caller("failed to read template directory", err,1)
		return err
	}
	for _, entry := range entries{
		template := strings.Split(entry.Name(), ".")[0]
		// 读取配置文件,读取错误则panic退出该进程
		viper.SetConfigFile(filepath.Join(project_dir.(string),"template",template + ".template.yaml"))
		err = viper.ReadInConfig()
		if err != nil {
			Logger_caller(fmt.Sprintf("Read %s failed!", template), err,1)
			fmt.Fprintf(os.Stderr, "Critical! Load the %s template has failed,exiting.",template)
			os.Exit(2)
		}
		Set_value(viper.AllSettings(),template)
	}
	return nil
}