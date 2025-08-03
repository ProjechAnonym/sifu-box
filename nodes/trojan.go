package nodes

import (
	"net/url"
)

func trojanFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	transport := make(map[string]any)
	tls := make(map[string]any)
	for k, v := range content {
		switch k {
		case "port":
			outbound["server_port"] = v
		case "password":
			outbound["password"] = v
		case "name":
			outbound["tag"] = v
		case "skip-cert-verify":
			tls["insecure"] = v
			tls["enabled"] = true
		case "sni":
			tls["server_name"] = v
			tls["enabled"] = true
		case "network":
			transport["type"] = v
		case "ws-opts":
			if opts, ok := v.(map[string]any); ok {
				if headers, ok := opts["headers"].(map[string]any); ok {
					if host, ok := headers["Host"].(string); ok {
						transport["host"] = host
					}
				}
				if path, ok := opts["path"].(string); ok {
					transport["path"] = path
				}
			}
		case "client-fingerprint":
		case "servername":
		case "tfo":
		default:
			outbound[k] = v
		}
	}
	outbound["tls"] = tls
	outbound["transport"] = transport
	return outbound
}
func trojanFromBase64(content *url.URL) map[string]any {
	outbound := make(map[string]any)
	transport := make(map[string]any)
	tls := make(map[string]any)
	tls["enabled"] = true
	tls["server_name"] = content.Query().Get("sni")
	if content.Query().Get("allowInsecure") != "1" {
		tls["insecure"] = false
	} else {
		tls["insecure"] = true
	}
	outbound["tag"] = content.Fragment
	outbound["server"] = content.Hostname()
	outbound["server_port"] = content.Port()
	outbound["password"] = content.User.String()
	outbound["type"] = "trojan"
	if content.Query().Get("type") != "" {
		transport["type"] = content.Query().Get("type")
		transport["host"] = content.Query().Get("host")
		transport["path"] = content.Query().Get("path")
		outbound["transport"] = transport
	}
	outbound["tls"] = tls
	return outbound
}
