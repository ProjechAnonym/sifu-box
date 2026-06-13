package nodes

import (
	"encoding/base64"
	"fmt"
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
		case "plugin":
			if v, ok := v.(string); ok {
				switch v {
				case "obfs":
					outbound["plugin"] = "obfs-local"
				case "v2ray":
					outbound["plugin"] = "v2ray-plugin"
				}
			}
		case "plugin-opts":
			if opts, ok := v.(map[string]any); ok {
				pluginOpts := ""
				if mode, ok := opts["mode"].(string); ok {
					pluginOpts = fmt.Sprintf("obfs=%s;", mode)
				}
				if host, ok := opts["host"].(string); ok {
					pluginOpts += fmt.Sprintf("obfs-host=%s;", host)
				}
				if pluginOpts != "" {
					outbound["plugin_opts"] = pluginOpts
				}
			}
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
