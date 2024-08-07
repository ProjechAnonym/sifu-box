package singbox

import (
	"fmt"
	"sifu-box/utils"
	"strings"
)

// ParseYaml 解析给定的YAML内容,从中提取代理配置信息
// 参数:
//   content - 一个接口切片,包含了YAML解析后的配置内容
//   name - 一个字符串,理论上应该用于标识配置的名称,但在函数内部未使用
// 返回值:
//   一个映射切片,包含了解析后的代理配置
//   如果解析过程中遇到错误,会返回具体的错误信息
func ParseYaml(content []interface{}, name string) ([]map[string]interface{}, error) {
    // 检查content是否为空,如果为空,则返回错误
	if len(content) == 0 {
		return nil, fmt.Errorf("没有节点信息")
	}
    // 初始化一个空的映射切片,用于存储解析后的代理配置
	var proxies []map[string]interface{}
    // 遍历content中的每个代理配置
	for _, proxy := range content {
        // 将当前代理配置从接口类型断言为映射,并格式化YAML内容
        result, err := formatYaml(proxy.(map[string]interface{}))
        // 如果没有错误发生,则将格式化后的代理配置添加到proxies切片中
        if err == nil {
            proxies = append(proxies, result)
        }
    }
    // 返回解析后的代理配置切片,如果没有错误发生,则返回nil
	return proxies, nil
}

// formatYaml 根据提供的映射格式化并返回相应的代理配置
// 该函数处理yaml格式的代理配置,并根据配置类型（如vmess、ss、trojan）进行相应的解析和格式化
// 参数proxyMap是一个包含代理配置的映射,其中类型(type)是必选项,用于确定代理类型
// 返回值proxy是格式化后的代理配置映射,如果无法解析或发生错误,则返回错误信息err
func formatYaml(proxyMap map[string]interface{}) (proxy map[string]interface{},err error) {
    // 使用defer和recover捕获并处理函数内部可能发生的panic,确保函数能够安全退出
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
            utils.LoggerCaller("Panic occurred in FormatUrl", err, 1)
            proxy = nil
            return
        }
    }()

    // 获取代理类型,这是确定如何解析代理配置的关键
    protocolType := proxyMap["type"]

    // 根据代理类型执行相应的解析逻辑
    switch protocolType {
    case "vmess":
        // 解析vmess类型的代理配置
        vmess,err := MarshalVmess(proxyMap)
        if err != nil{
            utils.LoggerCaller("解析vmess失败",err,1)
            return nil,err
        }
        proxy = vmess
    case "ss":
        // 解析ss（shadowsocks）类型的代理配置
        ss,err := MarshalShadowsocks(proxyMap)
        if err != nil{
            utils.LoggerCaller("解析shadowsocks失败",err,1)
            return nil,err
        }
        proxy = ss
    case "trojan":
        // 解析trojan类型的代理配置
        trojan,err := MarshalTrojan(proxyMap)
        if err != nil{
            utils.LoggerCaller("解析trojan失败",err,1)
            return nil,err
        }
        proxy = trojan
    default:
        // 如果代理类型未预置,则记录错误并返回
        utils.LoggerCaller("协议未预置",fmt.Errorf("没有预置'%s'协议", protocolType),1)
        return nil, fmt.Errorf("没有预置'%s'协议", protocolType)
    }
    // 返回格式化后的代理配置和可能的错误信息
    return proxy, err
}
// ParseUrl 解析URL列表,并返回解析后的代理配置
// 参数:
//   urls - 一个URL字符串的切片,代表待解析的URL列表
//   name - 一个字符串,代表名称,当前函数实现中未使用
// 返回值:
//   一个切片,其中包含解析后的URLs,每个URL对应一个map[string]interface{}
//   如果发生错误（例如,URL列表为空）,则返回nil和错误信息
func ParseUrl(urls []string, name string) ([]map[string]interface{}, error) {
    // 检查URL列表是否为空
    if len(urls) == 0 {
        return nil, fmt.Errorf("没有节点信息")
    }
    // 初始化一个切片,用于存储解析后的代理配置
    var proxies []map[string]interface{}

    // 遍历URL列表,对每个URL进行解析
    for _, url := range urls {
        // 尝试格式化URL,并获取解析结果
        result, err := formatUrl(url)
        // 如果没有发生错误,则将解析结果添加到proxies切片中
        if err == nil {
            proxies = append(proxies, result)
        }
        // 注意：这里没有处理err!=nil的情况,意味着如果有错误URL,会默默失败
    }
    // 返回解析后的代理配置和nil错误
    return proxies, nil
}
// formatUrl 根据提供的URL格式化代理设置
// 它支持ss、vmess和trojan协议
// 参数url是待格式化的URL字符串
// 返回值proxy是一个映射,包含解析后的代理配置
// 如果发生错误,返回值err将包含错误信息
func formatUrl(url string)(proxy map[string]interface{},err error){
    // 捕获并处理函数内的panic,转换为error返回
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
            utils.LoggerCaller("Panic occurred in FormatUrl", err, 1)
            proxy = nil
            return
        }
    }()
    
    // 提取URL中的协议类型
    protocolType := strings.Split(url, "://")[0]
    
    // 根据协议类型执行相应的解析操作
    switch protocolType {
    case "ss":
        // 解析shadowsocks链接
        ss,err := Base64Shadowsocks(url)
        if err != nil{
            utils.LoggerCaller("解析shadowsocks失败",err,1)
            return nil,err
        }
        proxy = ss
    case "vmess":
        // 解析vmess链接
        vmess,err := Base64Vmess(url)
        if err != nil{
            utils.LoggerCaller("解析vmess失败",err,1)
            return nil,err
        }
        proxy = vmess
    case "trojan":
        // 解析trojan链接
        trojan,err := Base64Trojan(url)
        if err != nil {
            utils.LoggerCaller("解析trojan失败",err,1)
            return nil,err
        }
        proxy = trojan
    default:
        // 如果协议类型不支持,则记录错误并返回
        utils.LoggerCaller("协议未预置",fmt.Errorf("没有预置'%s'协议", protocolType),1)
        return nil, fmt.Errorf("没有预置'%s'协议", protocolType)
    }
    return proxy, err
    
}