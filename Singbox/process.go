package singbox

import (
	"encoding/json"
	"fmt"
	utils "sifu-box/Utils"
	"strings"

	"github.com/huandu/go-clone"
)

// Get_map_value 根据提供的键序列从嵌套映射中检索值
// 它接受一个映射和一个或多个键,然后尝试按照键的顺序深入映射获取最终值
// 如果某个键不存在,函数将返回一个错误
// 参数:
//   proxy_map: 嵌套映射的初始副本
//   keys: 用于深入映射的键序列
// 返回值:
//   interface{}: 成功获取到的最终值,如果未成功,则为nil
//   error: 如果在查找过程中遇到任何键不存在的情况,则返回错误信息
func Get_map_value(proxy_map map[string]interface{},keys ...string)(interface{},error){
	// 使用clone库克隆初始映射,以避免对原始映射的意外修改
	result := clone.Clone(proxy_map)
	// 遍历提供的键序列
	for i, key := range keys {
		// 尝试获取当前键对应的值,并更新result变量
		// 如果值不存在,result将为nil,此时返回错误
		if result = result.(map[string]interface{})[key]; result == nil{
			return nil,fmt.Errorf("the key %s for level %d does not exist",key,i + 1)
		}
	}
	// 如果所有键都成功找到,返回最终值
	return result, nil
}
func Struct2map[T trojan|vmess|shadowsocks](s T,class string) (map[string]interface{},error){
	// 将s配置结构体序列化为JSON格式
	s_bytes, err := json.Marshal(s)
	if err != nil{
		utils.Logger_caller(fmt.Sprintf("marshal %s struct failed",class),err,1)
		return nil,err
	}
	// 反序列化JSON数据回map[string]interface{}格式,以便于后续处理
	var s_map map[string]interface{}
	
	err = json.Unmarshal(s_bytes, &s_map)
	if err != nil {
		utils.Logger_caller(fmt.Sprintf("marshal %s struct to map failed",class),err,1)
		return nil,err
	}
	return s_map,nil
}
// Format_yaml 根据给定的代理配置映射和模板字符串,格式化并返回相应的代理配置
// proxy_map: 包含代理配置信息的映射,如协议类型、服务器地址等
// template: 包含代理配置模板的字符串
// 返回值:
// - 格式化后的代理配置映射
// - 如果格式化失败,返回错误信息
func Format_yaml(proxy_map map[string]interface{},template string) (proxy map[string]interface{},err error) {
	// 使用defer和recover处理函数内部可能出现的panic,确保函数能够安全返回
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			utils.Logger_caller("Panic occurred in FormatUrl", err, 1)
			proxy = nil
			return
		}
	}()
	// 从proxy_map中提取协议类型和标签信息
	// 获取协议类型
	protocol_type := proxy_map["type"]

	// 根据协议类型切换不同的处理逻辑
	switch protocol_type {
	case "vmess":
		vmess,err := Map_marshal_vmess(proxy_map)
		if err != nil{
			return nil,err
		}
		proxy = vmess
	case "ss":
		ss,err := Map_marshal_ss(proxy_map)
		if err != nil{
			return nil,err
		}
		proxy = ss
	case "trojan":
		trojan,err := Map_marshal_trojan(proxy_map)
		if err != nil{
			return nil,err
		}
		proxy = trojan
	default:
		// 如果协议类型不在支持的范围内,返回错误
		return nil, fmt.Errorf("protocol %s is not in the template", protocol_type)
	}
	return proxy, err
}

// Format_url 根据给定的链接和模板,解析链接并返回符合模板格式的配置信息
// link: 需要解析的链接
// template: 配置模板,用于生成最终的配置信息
// 返回值:
//   - 一个map[string]interface{},包含解析后的配置信息
//   - 一个error,如果解析过程中出现错误,则返回错误信息
func Format_url(link string, template string) (proxy map[string]interface{},err error) {
	// 使用defer和recover处理函数内部可能出现的panic,确保函数能够安全返回
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			utils.Logger_caller("Panic occurred in FormatUrl", err, 1)
			proxy = nil
			return
		}
	}()
	// 解析链接的协议类型
	protocol_type := strings.Split(link, "://")[0]
	switch protocol_type {
	case "ss":
		ss,err := Base64_marshal_ss(link)
		if err != nil{
			return nil,err
		}
		proxy = ss
	case "vmess":
		vmess,err := Base64_marshal_vmess(link)
		if err != nil{
			return nil,err
		}
		proxy = vmess
	case "trojan":
		trojan,err := Base64_marshal_trojan(link)
		if err != nil {
			return nil,err
		}
		proxy = trojan
	default:
		// 如果协议类型不在支持的范围内,则返回错误
		return nil, fmt.Errorf("protocol %s is not in the template", protocol_type)
	}
	return proxy, err
}