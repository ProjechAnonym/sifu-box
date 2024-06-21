package singbox

import (
	"fmt"
	utils "sifu-box/Utils"
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
	tls_enable,err := Get_map_value(proxy_map,"skip-cert-verify")
	if err != nil{
		utils.Logger_caller(fmt.Sprintf("%s has no tls key",proxy_map["name"]),err,1)
		tls_enable = false
	}
	// 再次获取skip-cert-verify配置,这里重复是为了提供更清晰的错误信息
	skip_cert_verify,err := Get_map_value(proxy_map,"skip-cert-verify")
	if err != nil{
		utils.Logger_caller(fmt.Sprintf("%s has no skip_cert_verify key",proxy_map["name"]),err,1)
		skip_cert_verify = false
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
	vmess_map,err := Struct2map(vmess,"vmess")
	if err != nil {
		utils.Logger_caller("marshal vmess to map failed",err,1)
		return nil,err
	}
	return vmess_map,nil
}