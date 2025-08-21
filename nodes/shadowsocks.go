package nodes

import (
	"encoding/base64"
	"net/url"
	"strconv"
	"strings"
)

func shadowsocksFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	for k, v := range content {
		switch k {
		case "port":
			if _, ok := v.(int); !ok {
				outbound["server_port"] = 0
				continue
			}
			outbound["server_port"] = v.(int)
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

func shadowsocksFromBase64(content *url.URL) map[string]any {
	outbound := make(map[string]any)
	outbound["tag"] = content.Fragment
	outbound["server"] = content.Hostname()
	port, err := strconv.Atoi(content.Port())
	if err != nil {
		return nil
	}

	outbound["server_port"] = port
	message, err := base64.RawURLEncoding.DecodeString(content.User.String())
	if err != nil {
		return nil
	}
	if len(strings.Split(string(message), ":")) < 2 {
		return nil
	}
	outbound["method"] = strings.Split(string(message), ":")[0]
	outbound["password"] = strings.Split(string(message), ":")[1]
	outbound["type"] = "shadowsocks"
	return outbound
}
