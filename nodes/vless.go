package nodes

import (
	"net/url"
)

func vlessFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	transport := make(map[string]any)
	for k, v := range content {
		switch k {
		case "port":
			outbound["server_port"] = v
		case "name":
			outbound["tag"] = v
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

func vlessFromBase64(content *url.URL) map[string]any {

	outbound := make(map[string]any)
	transport := make(map[string]any)
	outbound["type"] = "vless"
	outbound["server"] = content.Hostname()
	outbound["server_port"] = content.Port()
	outbound["tag"] = content.Fragment
	outbound["uuid"] = content.User.String()
	transport["type"] = content.Query().Get("type")
	transport["host"] = content.Query().Get("host")
	transport["path"] = content.Query().Get("path")
	outbound["transport"] = transport
	return outbound
}
