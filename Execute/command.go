package execute

import (
	"fmt"

	utils "sifu-box/Utils"
	"strings"
)

// Reload_config 重新加载sing-box服务的配置
// 如果加载失败或日志中包含错误信息,则记录错误并返回错误
// 返回值:
//   error: 如果重新加载配置失败或存在错误日志,则返回相应的错误;否则返回nil
func Reload_config(service string,server utils.Server) (bool,error){
	status := true
	var results,errors []string
	if server.Localhost{
		// 使用systemctl命令重新加载sing-box服务
		_,_,err := utils.Command_exec("systemctl", "reload", service)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("reload config failed!",err,1)
			return false,err
		}

		// 使用journalctl命令获取sing-box服务最近10条日志
		results,errors,err = utils.Command_exec("journalctl", "-u", service,"-n","1")
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("get journal log failed!",err,1)
			return false,err
		}
	}else{
		_,_,err := utils.Command_ssh(server,"systemctl","reload",service)
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("reload config failed!",err,1)
			return false,err
		}
		results,errors,err = utils.Command_ssh(server,"journalctl","-u",service,"-n","1")
		if err != nil {
			// 如果命令执行失败,则记录错误并返回
			utils.Logger_caller("get journal log failed!",err,1)
			return false,err
		}
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
	
    // 如果一切正常,则返回nil
    return status,nil
}

// Boot_service 启动指定的服务
// 如果服务器是本地主机,则直接使用systemctl启动服务；
// 如果服务器是远程主机,则通过SSH方式启动服务
// 参数:
//   service string - 需要启动的服务名称
//   server utils.Server - 服务器配置信息,包括是否为本地主机等
// 返回值:
//   error - 如果启动服务失败,则返回错误信息；否则返回nil
func Boot_service(service string,server utils.Server) error{
	// 初始化状态变量和错误变量
	var status bool
	var err error

	// 根据服务器是否为本地主机,选择合适的启动服务方法
	if server.Localhost{
		// 使用systemctl命令在本地主机启动服务
		_,_,err = utils.Command_exec("systemctl", "start", service)
	}else{
		// 通过SSH方式在远程主机启动服务
		_,_,err = utils.Command_ssh(server,"systemctl","start",service)
	}

	// 检查启动服务过程中是否发生错误
	if err != nil {
		// 如果发生错误,则记录错误信息并返回错误
		utils.Logger_caller("reload config failed!",err,1)
		return err
	}

	// 检查服务是否成功启动
	status,err = Check_service(service,server)
	if err != nil {
		// 如果检查服务状态发生错误,则记录错误信息并返回错误
		utils.Logger_caller(fmt.Sprintf("%s is not running",service),err,1)
		return err
	}

	// 如果服务状态检查失败,则返回错误信息
	if !status {
		return fmt.Errorf("%s service is dead",service)
	}

	// 如果一切顺利,返回nil表示服务启动成功
	return nil
}

// Check_service 检查指定服务是否在给定的服务器上运行
// service: 需要检查的服务名称
// server: 服务器的信息,包括是否为本地服务器和其他连接信息
// 返回值: 一个布尔值表示服务是否运行,一个错误对象表示检查过程中是否发生错误
func Check_service(service string,server utils.Server) (bool,error){
	// 初始化服务运行状态为false
	status := false
	// 初始化用于存储命令执行结果和错误信息的切片
	var results,errors []string
	var err error
	// 根据服务器是否为本地服务器,选择不同的方式检查服务状态
	if server.Localhost{
		// 对于本地服务器,直接使用systemctl命令检查服务状态
		results,errors,err = utils.Command_exec("systemctl", "status", service)
	}else{
		// 对于远程服务器,通过SSH连接执行systemctl命令检查服务状态
		results,errors,err = utils.Command_ssh(server,"systemctl", "status", service)
	}
	// 检查命令执行过程中是否发生错误
	if err != nil {
		// 记录错误并返回
		utils.Logger_caller("check service failed!",err,1)
		return false,err
	}
	// 检查命令执行的错误输出是否为空
	if len(errors) != 0{
		// 记录错误并返回
		utils.Logger_caller("error",fmt.Errorf("pipe has output error msg"),1)
		return false,fmt.Errorf("pipe has output error msg")
	}
	// 遍历命令执行的结果,检查服务是否处于运行状态
	for _,result := range(results){
		// 如果结果中包含"active (running)",表示服务正在运行
		if strings.Contains(result,"active (running)"){
			status = true
			// 设置状态为运行并终止循环
			break
		}
	}
	// 返回服务运行状态和错误对象
	return status,nil
}