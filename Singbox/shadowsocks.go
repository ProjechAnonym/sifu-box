package singbox

import (
	"encoding/base64"
	"fmt"
	"net/url"
	utils "sifu-box/Utils"
	"strconv"
	"strings"
)

type shadowsocks struct {
	Type        string `json:"type"`
	Tag         string `json:"tag"`
	Server      string `json:"server"`
	Server_port int    `json:"server_port"`
	Method      string `json:"method"`
	Password    string `json:"password"`
}

// Map_marshal_ss 将一个包含代理配置信息的map转换为shadowsocks配置结构体,并返回该结构体的map表示形式
// 参数:
//   proxy_map: 一个包含代理配置信息的map,其中键名对应配置项的名称,键值为配置项的值
// 返回值:
//   转换后的shadowsocks配置结构体的map表示形式
//   如果转换过程中发生错误,返回错误信息
func Map_marshal_ss(proxy_map map[string]interface{}) (map[string]interface{}, error) {
    // 创建一个shadowsocks配置结构体实例,初始化其字段值从proxy_map中获取
    ss := shadowsocks{
        Type:        "shadowsocks",
        Tag:         proxy_map["name"].(string),
        Server:      proxy_map["server"].(string),
        Server_port: proxy_map["port"].(int),
        Method:      proxy_map["cipher"].(string),
        Password:    proxy_map["password"].(string),
    }

    // 将shadowsocks配置结构体转换为map,便于后续处理或返回
    // 这里使用了Struct2map函数进行转换,如果转换失败,则记录错误日志并返回错误
    ss_map, err := Struct2map(ss, "ss")
    if err != nil {
        utils.Logger_caller("marshal vmess to map failed", err, 1)
        return nil, err
    }

    // 转换成功,返回转换后的map以及nil错误
    return ss_map, nil
}

// Base64_marshal_ss 解析ss链接，并将其转换为map格式。
// link: ss链接。
// 返回值:
// - 一个包含解析后信息的map[string]interface{}。
// - 如果解析过程中出现错误，则返回错误信息。
func Base64_marshal_ss(link string) (map[string]interface{}, error) {
    // 移除链接前缀"ss://"并解码URL编码的部分
	info, err := url.QueryUnescape(strings.TrimPrefix(link, "ss://"))
	if err != nil {
        // 记录日志并返回错误
		utils.Logger_caller("url tag unescape failed", err, 1)
		return nil, err
	}

    // 根据"@"分割解码后的信息和服务器信息
	parts := strings.Split(info, "@")
    // 解码信息部分
	decoded_info, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
        // 记录日志并返回错误
		utils.Logger_caller("base64 decode failed", err, 1)
		return nil, err
	}

    // 分割解码后的信息为方法和密码
	info_parts := strings.Split(string(decoded_info), ":")
    // 分割服务器信息为地址和标签
	server_info := strings.Split(parts[1], "#")

    // 解析服务器URL
	server_url, err := url.Parse("ss://" + server_info[0])
	if err != nil {
        // 返回解析服务器URL失败的错误
		return nil, fmt.Errorf("failed to parse server URL: %v", err)
	}

    // 获取服务器端口
	port, err := strconv.Atoi(server_url.Port())
	if err != nil {
        // 记录日志并返回错误
		utils.Logger_caller("parse port failed", err, 1)
		return nil, err
	}

    // 构建shadowsocks配置结构体
	ss := shadowsocks{
		Type: "ss",
		Tag: server_info[1],
		Server: server_url.Hostname(),
		Server_port: port,
		Method: info_parts[0],
		Password: info_parts[1],
	}

    // 将shadowsocks结构体转换为map
	ss_map, err := Struct2map(ss, "shadowsocks")
	if err != nil {
        // 记录日志并返回错误
		utils.Logger_caller("marshal vmess to map failed", err, 1)
		return nil, err
	}

    // 返回转换后的map和nil错误
	return ss_map, nil
}