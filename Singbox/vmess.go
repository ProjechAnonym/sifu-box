package singbox

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	utils "sifu-box/Utils"
	"strconv"
	"strings"
)

type vmess struct {
	Type        string `json:"type"`
	Tag         string `json:"tag"`
	Server      string `json:"server"`
	Server_port int `json:"server_port"`
	Uuid        string `json:"uuid"`
	Alter_id  	int `json:"alter_id"`
	Security 	string `json:"security"`
	Tls 		Tls_config `json:"tls"`
	Transport   interface{} `json:"transport"`
}


// Map_marshal_vmess 根据给定的映射配置生成vmess协议配置的JSON格式
// proxy_map: 包含vmess代理配置的映射
// 返回值:
//   - map[string]interface{}: 生成的vmess配置的映射表示
//   - error: 如果配置生成过程中出现错误,则返回错误信息
func Map_marshal_vmess(proxy_map map[string]interface{}) (map[string]interface{},error) {
    // 尝试获取是否启用TLS的配置,如果不存在则默认为false
	tls_enable,err := Get_map_value(proxy_map,"tls")
	if err != nil{
		utils.Logger_caller(fmt.Sprintf("%s has no tls key",proxy_map["name"]),err,1)
		tls_enable = false
	}
	// 获取skip-cert-verify配置
	skip_cert_verify,err := Get_map_value(proxy_map,"skip-cert-verify")
	if err != nil{
		utils.Logger_caller(fmt.Sprintf("%s has no skip_cert_verify key",proxy_map["name"]),err,1)
		skip_cert_verify = true
	}
	// 初始化vmess配置结构体
	vmess := vmess{
		Tag:         proxy_map["name"].(string),
		Type:        "vmess",
		Server:      proxy_map["server"].(string),
		Server_port: proxy_map["port"].(int),
		Uuid:        proxy_map["uuid"].(string),
		Alter_id: 	 proxy_map["alterId"].(int),
		Security: 	 proxy_map["cipher"].(string),
		Tls: Tls_config{
			Enabled:  tls_enable.(bool),
			Insecure: skip_cert_verify.(bool),
		},
	}
	// 根据网络类型配置传输方式
	switch proxy_map["network"].(string) {
	case "grpc":
		// 配置gRPC传输选项
		service_name,err := Get_map_value(proxy_map,"grpc-opts","grpc-service-name")
		if err != nil{
			return nil,err
		}
		transport := Grpc{
			Type:                  proxy_map["network"].(string),
			Service_name:          service_name.(string),
			Idle_timeout:          "15s",
			Ping_timeout:          "15s",
			Permit_without_stream: false,
		}
		vmess.Transport = transport
	case "ws":
		// 配置WebSocket传输选项
		transport := Websocket{
			Type:                   proxy_map["network"].(string),
			Path:                   proxy_map["ws-path"].(string),
			Headers:                map[string]string{"host": proxy_map["ws-headers"].(map[string]interface{})["Host"].(string)},
			Early_data_header_name: "Sec-WebSocket-Protocol",
		}
		vmess.Transport = transport
	}
	// 转换为map字典
	vmess_map,err := Struct2map(vmess,"vmess")
	if err != nil {
		utils.Logger_caller("marshal vmess to map failed",err,1)
		return nil,err
	}
	return vmess_map,nil
}

// Base64_marshal_vmess 解析vmess协议链接,并将其转换为map格式,方便后续处理
// link: vmess协议链接,格式为"vmess://<base64编码的配置信息>"
// 返回值:
//   - map[string]interface{}: 包含vmess协议配置信息的map
//   - error: 解析过程中可能出现的错误
func Base64_marshal_vmess(link string) (map[string]interface{},error){
    // 移除链接前缀"vmess://",获取base64编码的配置信息
	info := strings.TrimPrefix(link, "vmess://")
    // 解码base64编码的配置信息
	var decoded_info []byte
	var err error
	decoded_info, err = base64.URLEncoding.DecodeString(info)
	if err != nil {
        // 记录解码失败的日志,并返回错误
		utils.Logger_caller("base64 decode failed",err,1)
		return nil,err
	}
    // 将解码后的信息反序列化为map格式
	var proxy_map map[string]interface{}
	if err := json.Unmarshal(decoded_info,&proxy_map);err != nil{
        // 记录反序列化失败的日志,并返回错误
		utils.Logger_caller("string convert to map failed",err,1)
		return nil,err
	}
    // 从map中提取并转换端口号和alter_id
	port,err := strconv.Atoi(proxy_map["port"].(string))
	if err != nil {
        // 记录转换失败的日志,并返回错误
		utils.Logger_caller("convert to num failed!",err,1)
		return nil,err
	}
	alter_id,err := strconv.Atoi(proxy_map["aid"].(string))
	if err != nil {
        // 记录转换失败的日志,并返回错误
		utils.Logger_caller("convert to num failed!",err,1)
		return nil,err
	}
    // 判断tls是否启用,以及是否需要跳过证书验证
	var tls_enable bool
	if _,err := Get_map_value(proxy_map,"tls"); err != nil {
        // 记录未找到tls键的日志,设置tls为禁用,并返回错误
		utils.Logger_caller(fmt.Sprintf("%s has no tls key",proxy_map["ps"].(string)),err,1)
		tls_enable = false
	}else{
		tls_enable = true
	}
	
	skip_cert,err := Get_map_value(proxy_map,"skip-cert-verify")
	if err != nil {
        // 记录未找到skip_cert_verify键的日志,设置为跳过证书验证,并返回错误
		utils.Logger_caller(fmt.Sprintf("%s has no skip_cert_verify key",proxy_map["ps"].(string)),err,1)
		skip_cert = true
	}
	sni,err := Get_map_value(proxy_map,"sni")
	if err != nil{
        // 记录未找到sni键的日志,设置sni为空字符串,并返回错误
		utils.Logger_caller(fmt.Sprintf("%s has no sni key",proxy_map["ps"].(string)),err,1)
		sni = ""
	}
    // 根据提取的信息,构建tls配置
	tls := Tls_config{
		Enabled: tls_enable,
		Insecure: skip_cert.(bool),
		Server_name: sni.(string),
	}
    // 根据提取的信息,构建vmess配置
	vmess := vmess{
		Type: "vmess",
		Tag: proxy_map["ps"].(string),
		Server: proxy_map["add"].(string),
		Server_port: port,
		Uuid: proxy_map["id"].(string),
		Alter_id: alter_id,
		Security: "auto",
		Tls: tls,
	}

    // 根据网络类型设置传输协议
	switch proxy_map["net"].(string) {
		case "ws":
			transport := Websocket{
				Type: proxy_map["net"].(string),
				Path: proxy_map["path"].(string),
				Headers: map[string]string{"host":proxy_map["host"].(string)},
				Early_data_header_name: "Sec-WebSocket-Protocol",
			}
			vmess.Transport = transport
		}
    // 将vmess配置转换为map格式,并返回
	vmess_map,err := Struct2map(vmess,"vmess")
	if err != nil {
        // 记录转换失败的日志,并返回错误
		utils.Logger_caller("marshal vmess to map failed",err,1)
		return nil,err
	}
	return vmess_map,nil
}