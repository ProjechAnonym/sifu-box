package nodes

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
)

func vmessFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	transport := make(map[string]any)
	for k, v := range content {
		switch k {
		case "port":
			if _, ok := v.(int); !ok {
				outbound["server_port"] = 0
				continue
			}
			outbound["server_port"] = v.(int)
		case "cipher":
			outbound["security"] = v
		case "name":
			outbound["tag"] = v
		case "alterId":
			if _, ok := v.(int); !ok {
				outbound["alter_id"] = 0
				continue
			}
			outbound["alter_id"] = v.(int)
		case "network":
			transport["type"] = v
		case "ws-opts":
			if opts, ok := v.(map[string]any); ok {
				if headers, ok := opts["headers"].(map[string]any); ok {
					if host, ok := headers["Host"].(string); ok {
						transport["headers"] = map[string]any{"host": host}
					}
				}
				if path, ok := opts["path"].(string); ok {
					transport["path"] = path
				}
			}
		case "client-fingerprint":
		case "tfo":
		case "tls":
		case "skip-cert-verify":
		case "servername":
		default:
			outbound[k] = v
		}
	}
	outbound["transport"] = transport
	return outbound

}
func vmessFromBase64(content *url.URL) (map[string]any, error) {
	data, err := base64.StdEncoding.DecodeString(content.Host)
	if err != nil {
		return nil, err
	}
	config := make(map[string]any)
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	outbound := make(map[string]any)
	transport := make(map[string]any)
	for k, v := range config {
		switch k {
		case "ps":
			outbound["tag"] = v
		case "add":
			outbound["server"] = v
		case "port":
			if _, ok := v.(int); !ok {
				outbound["server_port"] = 0
				continue
			}
			outbound["server_port"] = v.(int)
		case "aid":
			if _, ok := v.(int); !ok {
				outbound["alter_id"] = 0
				continue
			}
			outbound["alter_id"] = v.(int)

		case "id":
			outbound["uuid"] = v
		case "net":
			transport["type"] = v
		case "path":
			transport["path"] = v
		case "host":
			transport["host"] = v
		}
	}
	if transport["type"] == "ws" {
		transport["headers"] = map[string]any{"host": transport["host"]}
		delete(transport, "host")
	}
	outbound["type"] = "vmess"
	outbound["security"] = "auto" // 默认安全性为 auto
	if transport["type"] != "ws" && transport["type"] != "http" && transport["type"] != "quic" && transport["type"] != "grpc" && transport["type"] != "httpupgrade" {
		return outbound, nil
	}
	outbound["transport"] = transport
	return outbound, nil
}
