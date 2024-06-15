package singbox

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	utils "sifu-box/Utils"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

// Format_yaml 根据给定的代理配置映射和模板字符串,格式化并返回相应的代理配置
// proxy_map: 包含代理配置信息的映射,如协议类型、服务器地址等
// template: 包含代理配置模板的字符串
// 返回值:
// - 格式化后的代理配置映射
// - 如果格式化失败,返回错误信息
func Format_yaml(proxy_map map[string]interface{},template string) (map[string]interface{}, error) {
	// 从proxy_map中提取协议类型和标签信息
	// 获取协议类型
	protocol_type := proxy_map["type"]
	tag := proxy_map["name"]

	// 根据协议类型切换不同的处理逻辑
	switch protocol_type {
	case "vmess":
		// 处理vmess协议类型的代理配置
		proxy_vmess, err := utils.Get_value(template, "outbounds", "vmess")
		if err != nil {
			// 如果获取vmess模板失败,记录日志并返回错误
			utils.Logger_caller("Get vmess Template failed!", err,1)
			return nil, err
		}
		// 根据proxy_map中的信息填充vmess配置
		proxy_vmess.(map[string]interface{})["tag"] = tag
		proxy_vmess.(map[string]interface{})["server"] = proxy_map["server"]
		proxy_vmess.(map[string]interface{})["server_port"] = int(proxy_map["port"].(int))
		proxy_vmess.(map[string]interface{})["uuid"] = proxy_map["uuid"]
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["type"] = proxy_map["network"]
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["path"] = proxy_map["ws-path"]
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["headers"] = map[string]string{"host":proxy_map["ws-headers"].(map[string]interface{})["Host"].(string)}
		// 返回填充后的vmess配置
		return proxy_vmess.(map[string]interface{}), nil

	case "ss":
		// 处理ss协议类型的代理配置
		proxy_ss, err := utils.Get_value(template, "outbounds", "ss")
		if err != nil {
			// 如果获取ss模板失败,记录日志并返回错误
			utils.Logger_caller("Get ss Template failed!", err,1)
			return nil, err
		}
		// 根据proxy_map中的信息填充ss配置
		proxy_ss.(map[string]interface{})["tag"] = tag
		proxy_ss.(map[string]interface{})["server"] = proxy_map["server"]
		proxy_ss.(map[string]interface{})["server_port"] = int(proxy_map["port"].(int))
		proxy_ss.(map[string]interface{})["method"] = proxy_map["cipher"]
		proxy_ss.(map[string]interface{})["password"] = proxy_map["password"]
		// 返回填充后的ss配置
		return proxy_ss.(map[string]interface{}), nil

	case "trojan":
		// 处理trojan协议类型的代理配置
		proxy_trojan, err := utils.Get_value(template, "outbounds", "trojan")
		if err != nil {
			// 如果获取trojan模板失败,记录日志并返回错误
			utils.Logger_caller("Get trojan Template failed!", err,1)
			return nil, err
		}
		// 根据proxy_map中的信息填充trojan配置
		proxy_trojan.(map[string]interface{})["tag"] = tag
		proxy_trojan.(map[string]interface{})["server"] = proxy_map["server"]
		proxy_trojan.(map[string]interface{})["server_port"] = int(proxy_map["port"].(int))
		proxy_trojan.(map[string]interface{})["tls"].(map[string]interface{})["server_name"] = proxy_map["sni"]
		proxy_trojan.(map[string]interface{})["password"] = proxy_map["password"]
		// 返回填充后的trojan配置
		return proxy_trojan.(map[string]interface{}), nil
	}

	// 如果协议类型不在支持的范围内,返回错误
	return nil, fmt.Errorf("protocol %s is not in the template", protocol_type)
}

// Format_url 根据给定的链接和模板,解析链接并返回符合模板格式的配置信息
// link: 需要解析的链接
// template: 配置模板,用于生成最终的配置信息
// 返回值:
//   - 一个map[string]interface{},包含解析后的配置信息
//   - 一个error,如果解析过程中出现错误,则返回错误信息
func Format_url(link string, template string) (map[string]interface{}, error) {
    // 使用defer和recover处理函数内部可能出现的panic,确保函数能够安全返回
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered from panic: %v", r)
			utils.Logger_caller("Panic occurred in FormatUrl", err, 1)
		}
	}()

	// 解析链接的协议类型
	protocol_type := strings.Split(link, "://")[0]
	switch protocol_type {
	case "ss":
		// 解析ss链接格式,并根据解析结果生成相应的配置信息
		re := regexp.MustCompile(`ss:\/\/([^@]+)@([^:]+):(\d+)(#.*)?`)
		matches := re.FindStringSubmatch(link)
		if matches == nil {
			utils.Logger_caller("split ss url failed", nil, 1)
			return nil, fmt.Errorf("link %s is not in the format of ss://", link)
		}
		tag, err := url.QueryUnescape(matches[4])
		if err != nil {
			utils.Logger_caller("url tag unescape failed", err, 1)
			return nil, err
		}
		msg_bytes, err := base64.StdEncoding.DecodeString(matches[1] + "=")
		if err != nil {
			utils.Logger_caller("url password and cipher unescape failed", err, 1)
			return nil, err
		}
		proxy_ss, err := utils.Get_value(template, "outbounds", "ss")
		if err != nil {
			utils.Logger_caller("Get ss Template failed!", err, 1)
			return nil, err
		}
		proxy_ss.(map[string]interface{})["tag"] = tag
		proxy_ss.(map[string]interface{})["server"] = matches[2]
		proxy_ss.(map[string]interface{})["server_port"], err = strconv.Atoi(matches[3])
		if err != nil {
			utils.Logger_caller("num string transfer failed!", err, 1)
			return nil, err
		}
		proxy_ss.(map[string]interface{})["method"] = strings.Split(string(msg_bytes), ":")[0]
		proxy_ss.(map[string]interface{})["password"] = strings.Split(string(msg_bytes), ":")[1]
		return proxy_ss.(map[string]interface{}), nil
	case "vmess":
		// 解析vmess链接格式,并根据解析结果生成相应的配置信息
		msg_bytes, err := base64.StdEncoding.DecodeString(strings.Split(link, "://")[1])
		if err != nil {
			utils.Logger_caller("url password and cipher unescape failed", err, 1)
			return nil, err
		}
		msg, err := simplejson.NewJson(msg_bytes)
		if err != nil {
			utils.Logger_caller("vmess msg unescape failed", err, 1)
			return nil, err
		}
		proxy_vmess, err := utils.Get_value(template, "outbounds", "vmess")
		if err != nil {
			utils.Logger_caller("Get vmess Template failed!", err, 1)
			return nil, err
		}
		proxy_vmess.(map[string]interface{})["tag"] = msg.Get("ps")
		proxy_vmess.(map[string]interface{})["server"] = msg.Get("add")
		proxy_vmess.(map[string]interface{})["server_port"] = msg.Get("port").MustInt()
		proxy_vmess.(map[string]interface{})["uuid"] = msg.Get("id")
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["type"] = msg.Get("net")
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["path"] = msg.Get("path")
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["headers"] = map[string]string{"host": msg.Get("host").MustString()}
		return proxy_vmess.(map[string]interface{}), nil
	case "trojan":
		// 解析trojan链接格式,并根据解析结果生成相应的配置信息
		re := regexp.MustCompile(`^(.*?)://([^@]+)@([^:]+):(\d+)\?(.*?)#(.*)$`)
		matches := re.FindStringSubmatch(link)
		tag, err := url.QueryUnescape(matches[6])
		if err != nil {
			utils.Logger_caller("url tag unescape failed", err, 1)
			return nil, err
		}
		proxy_trojan, err := utils.Get_value(template, "outbounds", "trojan")
		if err != nil {
			utils.Logger_caller("Get trojan Template failed!", err, 1)
			return nil, err
		}
		proxy_trojan.(map[string]interface{})["tag"] = tag
		proxy_trojan.(map[string]interface{})["server"] = matches[3]
		proxy_trojan.(map[string]interface{})["server_port"], err = strconv.Atoi(matches[4])
		if err != nil {
			utils.Logger_caller("num string transfer failed!", err, 1)
			return nil, err
		}
		values, err := url.ParseQuery(matches[5])
		if err != nil {
			utils.Logger_caller("sni string parse failed!", err, 1)
			return nil, err
		}
		sniValue := values.Get("sni")
		proxy_trojan.(map[string]interface{})["tls"].(map[string]interface{})["server_name"] = sniValue
		proxy_trojan.(map[string]interface{})["password"] = matches[2]
		return proxy_trojan.(map[string]interface{}), nil
	}

	// 如果协议类型不在支持的范围内,则返回错误
	return nil, fmt.Errorf("protocol %s is not in the template", protocol_type)
}