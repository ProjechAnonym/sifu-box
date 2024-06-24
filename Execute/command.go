package execute

import (
	"fmt"
	database "sifu-box/Database"
	utils "sifu-box/Utils"
	"strings"
)

// Reload_config 重新加载sing-box服务的配置
// 如果加载失败或日志中包含错误信息,则记录错误并返回错误
// 返回值:
//   error: 如果重新加载配置失败或存在错误日志,则返回相应的错误;否则返回nil
func Reload_config(service string,server database.Server) (bool,error){
	status := true
	if server.Localhost{
		// 使用systemctl命令重新加载sing-box服务
		_,_,err := utils.Command_exec("systemctl", "reload", service)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("reload config failed!",err,1)
			return false,err
		}

		// 使用journalctl命令获取sing-box服务最近10条日志
		results,errors,err := utils.Command_exec("journalctl", "-u", service,"-n","1")
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("get journal log failed!",err,1)
			return false,err
		}
		 
		// 检查日志中是否包含错误信息
		for _,result := range(results){
			if strings.Contains(result,"ERROR"){
				// 如果日志包含错误信息,则记录错误并返回
				utils.Logger_caller("reload config failed!",fmt.Errorf(result),1)
				status = false
				break
			}
		}

		// 如果命令的标准错误输出不为空,则记录错误并返回
		if len(errors) != 0{
			utils.Logger_caller("error",fmt.Errorf("pipe has output error msg"),1)
			return false,fmt.Errorf("pipe has output error msg")
		}

		// 如果没有错误但服务未成功重新加载,则返回相应错误
		if !status{
			return false,fmt.Errorf("reload new config failed")
		}
	}
    // 如果一切正常,则返回nil
    return status,nil
}

func Boot_service(service string,server database.Server) error{
	var status bool
	if server.Localhost{
		// 使用systemctl命令重新加载sing-box服务
		_,_,err := utils.Command_exec("systemctl", "start", service)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("reload config failed!",err,1)
			return err
		}
		status,err = Check_service(service,server)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller(fmt.Sprintf("%s is not running",service),err,1)
			return err
		}

	}
	if !status {
		return fmt.Errorf("%s service is dead",service)
	}
	return nil
}

func Check_service(service string,server database.Server) (bool,error){
	status := false
	if server.Localhost{
		// 使用systemctl命令重新加载sing-box服务
		results,errors,err := utils.Command_exec("systemctl", "status", service)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("check service failed!",err,1)
			return false,err
		}
		// 如果命令的标准错误输出不为空,则记录错误并返回
		if len(errors) != 0{
			utils.Logger_caller("error",fmt.Errorf("pipe has output error msg"),1)
			return false,fmt.Errorf("pipe has output error msg")
		}
		// 检查日志中是否包含运行信息
		for _,result := range(results){
			if strings.Contains(result,"active (running)"){
				status = true
				break
			}
		}
	}
	return status,nil
}