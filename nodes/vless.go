package nodes

import (
	"net/url"
	"strconv"
)

func vlessFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	transport := Transport{}
	for k, v := range content {
		switch k {
		case "port":
			if _, ok := v.(int); !ok {
				outbound["server_port"] = 0
				continue
			}
			outbound["server_port"] = v.(int)
		case "name":
			outbound["tag"] = v
		case "network":
			if network, ok := v.(string); ok {
				transport.Type = network
			}

		case "ws-opts":
			if opts, ok := v.(map[string]any); ok {
				if headers, ok := opts["headers"].(map[string]any); ok {
					if host, ok := headers["Host"].(string); ok {
						transport.Headers = map[string]string{"host": host}
					}
				}
				if path, ok := opts["path"].(string); ok {
					transport.Path = path
				}
			}
		case "client-fingerprint":
		case "tfo":
		case "tls":
		case "skip-cert-verify":
		case "servername":
		case "udp":
			continue
		default:
			outbound[k] = v
		}
	}
	if transport.Type != "ws" && transport.Type != "http" && transport.Type != "quic" && transport.Type != "grpc" && transport.Type != "httpupgrade" {
		return outbound
	}
	outbound["transport"] = transport
	return outbound
}

func vlessFromBase64(content *url.URL) map[string]any {

	outbound := make(map[string]any)
	transport := Transport{}
	outbound["type"] = "vless"
	outbound["server"] = content.Hostname()
	port, err := strconv.Atoi(content.Port())
	if err != nil {
		return nil
	}

	outbound["server_port"] = port
	outbound["tag"] = content.Fragment
	outbound["uuid"] = content.User.String()
	transport.Type = content.Query().Get("type")
	host := content.Query().Get("host")
	transport.Path = content.Query().Get("path")
	if transport.Type == "ws" {
		transport.Headers = map[string]string{"host": host}
	}
	if transport.Type != "ws" && transport.Type != "http" && transport.Type != "quic" && transport.Type != "grpc" && transport.Type != "httpupgrade" {
		return outbound
	}
	outbound["transport"] = transport
	return outbound
}
