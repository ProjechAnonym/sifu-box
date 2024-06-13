package main

import (
	"fmt"

	"github.com/huandu/go-clone"
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
	base_dir := "/root/sifu-box"
	return base_dir
}

