package singbox

import (
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
        skip_cert_verify = false
    }

    // 尝试获取sni配置,用于TLS连接时指定服务器名称
    sni, err := Get_map_value(proxy_map, "sni")
    if err != nil {
        // 如果获取失败,记录错误日志并返回错误
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
        return nil, err
    }

    // 返回转换后的trojan配置map和nil错误
    return trojan_map, nil
}

// Base64_marshal_trojan 解析trojan链接并将其转换为map格式,用于后续处理
// link: trojan协议的链接字符串,例如"trojan://password@server#tag?allowInsecure=1&sni=example.com:8080"
// 返回值:
//   - 一个包含解析后信息的map[string]interface{},其中信息已编码为Base64
//   - 如果解析过程中出现错误,则返回错误信息
func Base64_marshal_trojan(link string) (map[string]interface{}, error) {
    // 移除链接前缀"trojan://",以便后续处理
    info := strings.TrimPrefix(link, "trojan://")
    // 使用"@"分割链接字符串,获取密码和服务器信息
    parts := strings.Split(info, "@")
    // 从分割后的第一部分获取密码
    password := parts[0]
    // 使用"#"分割服务器信息,获取服务器URL和标签
    url_parts := strings.Split(parts[1], "#")
    // 解析服务器URL,为后续获取端口和参数做准备
    server_url, err := url.Parse("trojan://" + url_parts[0])
    if err != nil {
        // 日志记录URL解析失败
        utils.Logger_caller("trojan url parsed failed!", err, 1)
        return nil, err
    }
    // 解码标签信息,以便正确使用
    tag, err := url.QueryUnescape(url_parts[1])
    if err != nil {
        // 日志记录标签解码失败
        utils.Logger_caller("url decode failed!", err, 1)
        return nil, err
    }
    // 获取服务器端口
    port, err := strconv.Atoi(server_url.Port())
    if err != nil {
        // 日志记录端口获取失败
        utils.Logger_caller("get trojan port failed!", err, 1)
        return nil, err
    }
    
    // 从服务器URL中获取参数
    params := server_url.Query()
    // 初始化是否跳过证书验证的变量
    var skip_cert bool
    // 根据参数"allowInsecure"的值,确定是否跳过证书验证
    if skip_cert_verify := params.Get("allowInsecure"); skip_cert_verify != "" {
        if skip_cert_verify == "1" {
            skip_cert = true
        } else {
            skip_cert = false
        }
    } else {
        skip_cert = true
    }
    // 构建trojan配置结构体
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
    // 将trojan配置结构体转换为map格式
    trojan_map, err := Struct2map(trojan, "trojan")
    if err != nil {
        // 日志记录结构体转换失败
        return nil, err
    }
    // 返回转换后的map和nil错误
    return trojan_map, nil
}