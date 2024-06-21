package singbox

import (
	"fmt"
	"net/url"
	utils "sifu-box/Utils"
	"strconv"
	"strings"
)

type trojan struct {
	Type        string     `json:"type"`
	Tag         string     `json:"tag"`
	Server      string     `json:"server"`
	Server_port int        `json:"server_port"`
	Password    string     `json:"password"`
	Tls         Tls_config `json:"tls"`
}

// Map_marshal_trojan 将给定的map转换为trojan配置的map格式
// 参数proxy_map包含trojan代理的配置信息
// 返回转换后的trojan配置map和可能的错误
func Map_marshal_trojan(proxy_map map[string]interface{}) (map[string]interface{}, error) {
    // 尝试获取skip-cert-verify配置,用于确定是否忽略证书验证
    skip_cert_verify, err := Get_map_value(proxy_map, "skip-cert-verify")
    if err != nil {
        // 如果获取失败,记录错误日志并使用默认值false
        utils.Logger_caller(fmt.Sprintf("%s has no skip-cert-verify key", proxy_map["name"].(string)), err, 1)
        skip_cert_verify = false
    }

    // 尝试获取sni配置,用于TLS连接时指定服务器名称
    sni, err := Get_map_value(proxy_map, "sni")
    if err != nil {
        // 如果获取失败,记录错误日志并返回错误
        utils.Logger_caller(fmt.Sprintf("%s has no sni key", proxy_map["sni"].(string)), err, 1)
        return nil, err
    }

    // 创建trojan配置结构体,填充从proxy_map获取的配置信息
    trojan := trojan{
        Type:        "trojan",
        Tag:         proxy_map["name"].(string),
        Server:      proxy_map["server"].(string),
        Server_port: proxy_map["port"].(int),
        Password:    proxy_map["password"].(string),
        Tls: Tls_config{
            Enabled:     true,
            Insecure:    skip_cert_verify.(bool),
            Server_name: sni.(string),
        },
    }

    // 将trojan配置结构体转换为map格式
    trojan_map, err := Struct2map(trojan, "trojan")
    if err != nil {
        // 如果转换失败,记录错误日志并返回错误
        utils.Logger_caller("marshal trojan to map failed", err, 1)
        return nil, err
    }

    // 返回转换后的trojan配置map和nil错误
    return trojan_map, nil
}

func Base64_marshal_trojan(link string) (map[string]interface{}, error){
	info := strings.TrimPrefix(link,"trojan://")
	parts := strings.Split(info, "@")
	password := parts[0]
	url_parts := strings.Split(parts[1], "#")
	server_url,err := url.Parse("trojan://" + url_parts[0])
	if err != nil {
		utils.Logger_caller("trojan url parsed failed!",err,1)
		return nil,err
	}
	tag,err := url.QueryUnescape(url_parts[1])
	if err != nil {
		utils.Logger_caller("url decode failed!",err,1)
		return nil,err
	}
	port,err := strconv.Atoi(server_url.Port())
	if err != nil {
		utils.Logger_caller("get trojan port failed!",err,1)
		return nil,err
	}
	
	params := server_url.Query()
	var skip_cert bool
	if skip_cert_verify := params.Get("allowInsecure");skip_cert_verify != ""{
		if skip_cert_verify == "1"{
			skip_cert = true
		}else{
			skip_cert = false
		}
	}else{
		skip_cert = true
	}
	trojan := trojan{
		Type: "trojan",
		Tag: tag,
		Password: password,
		Server: server_url.Hostname(),
		Server_port: port,
		Tls: Tls_config{
			Enabled: true,
			Insecure: skip_cert,
			Server_name: params.Get("sni"),
		},
	}
	trojan_map,err := Struct2map(trojan,"trojan")
	if err != nil{
		utils.Logger_caller("marshal trojan to map failed",err,1)
		return nil,err
	}
	return trojan_map,nil
}