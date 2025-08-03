package nodes

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

func shadowsocksFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	for k, v := range content {
		switch k {
		case "port":
			outbound["server_port"] = v
		case "cipher":
			outbound["method"] = v
		case "name":
			outbound["tag"] = v
		case "type":
			outbound["type"] = "shadowsocks"
		case "udp":
		default:
			outbound[k] = v
		}
	}
	return outbound
}

func shadowsocksFromBase64(content *url.URL) (map[string]any, error) {
	outbound := make(map[string]any)
	outbound["tag"] = content.Fragment
	outbound["server"] = content.Hostname()
	outbound["server_port"] = content.Port()
	message, err := base64.RawURLEncoding.DecodeString(content.User.String())
	if err != nil {
		return nil, err
	}
	if len(strings.Split(string(message), ":")) < 2 {
		return nil, fmt.Errorf("shadowsocks解密出错, 未能获得加密方法和密钥")
	}
	outbound["method"] = strings.Split(string(message), ":")[0]
	outbound["password"] = strings.Split(string(message), ":")[1]
	outbound["type"] = "shadowsocks"
	return outbound, nil
}
