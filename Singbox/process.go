package singbox

import (
	"errors"
	"fmt"
	utils "sifu-box/Utils"
)

func Format_yaml(proxy_map map[string]interface{},template string) (map[string]interface{}, error) {
	// 获取协议类型
	protocol_type := proxy_map["type"]
	tag := proxy_map["name"]
	switch protocol_type {
	case "vmess":
		// 获取模板信息
		proxy_vmess, err := utils.Get_value(template, "outbounds", "vmess")
		if err != nil {
			utils.Logger_caller("Get vmess Template failed!", err,1)
			return nil, err
		}
		proxy_vmess.(map[string]interface{})["tag"] = tag
		proxy_vmess.(map[string]interface{})["server"] = proxy_map["server"]
		proxy_vmess.(map[string]interface{})["server_port"] = int(proxy_map["port"].(int))
		proxy_vmess.(map[string]interface{})["uuid"] = proxy_map["uuid"]
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["type"] = proxy_map["network"]
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["path"] = proxy_map["ws-path"]
		proxy_vmess.(map[string]interface{})["transport"].(map[string]interface{})["headers"] = proxy_map["ws-headers"]
		return proxy_vmess.(map[string]interface{}), nil
	case "ss":
		// 获取模板信息
		proxy_ss, err := utils.Get_value(template, "outbounds", "ss")
		if err != nil {
			utils.Logger_caller("Get ss Template failed!", err,1)
			return nil, err
		}
		proxy_ss.(map[string]interface{})["tag"] = tag
		proxy_ss.(map[string]interface{})["server"] = proxy_map["server"]
		proxy_ss.(map[string]interface{})["server_port"] = int(proxy_map["port"].(int))
		proxy_ss.(map[string]interface{})["method"] = proxy_map["cipher"]
		proxy_ss.(map[string]interface{})["password"] = proxy_map["password"]
		return proxy_ss.(map[string]interface{}), nil
	case "trojan":
		// 获取模板信息
		proxy_trojan, err := utils.Get_value(template, "outbounds", "trojan")
		if err != nil {
			utils.Logger_caller("Get trojan Template failed!", err,1)
			return nil, err
		}
		proxy_trojan.(map[string]interface{})["tag"] = tag
		proxy_trojan.(map[string]interface{})["server"] = proxy_map["server"]
		proxy_trojan.(map[string]interface{})["server_port"] = int(proxy_map["port"].(int))
		proxy_trojan.(map[string]interface{})["tls"].(map[string]interface{})["server_name"] = proxy_map["sni"]
		proxy_trojan.(map[string]interface{})["password"] = proxy_map["password"]
		return proxy_trojan.(map[string]interface{}), nil
	}
	msg := fmt.Sprintf("protocol %s is not in the template", protocol_type)
	err := errors.New(msg)
	return nil, err
}