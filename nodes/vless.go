package nodes

import (
	"net/url"
	"strconv"
)

func vlessFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	transport := Transport{}
	tls := TLS{}
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
					if host, ok := headers["Host"]; ok {
						transport.Headers = map[string]any{"host": host}
					}
				}
				if path, ok := opts["path"].(string); ok {
					transport.Path = path
				}
				if path, ok := opts["path"].([]string); ok {
					transport.Path = path[0]
				}
			}
		case "http-opts":
			if opts, ok := v.(map[string]any); ok {
				if headers, ok := opts["headers"].(map[string]any); ok {
					if host, ok := headers["Host"]; ok {
						transport.Headers = map[string]any{"host": host}
					}
				}
				if path, ok := opts["path"].(string); ok {
					transport.Path = path
				}
				if path, ok := opts["path"].([]string); ok {
					transport.Path = path[0]
				}
			}
		case "reality-opts":
			if opts, ok := v.(map[string]any); ok {
				tls.Reality.Enabled = true
				if public_key, ok := opts["public-key"].(string); ok {
					tls.Reality.PublicKey = public_key
				}
				if short_id, ok := opts["short-id"].(string); ok {
					tls.Reality.ShortId = short_id
				}
			}
		case "skip-cert-verify":
			tls.Enabled = true
			tls.Insecure = true
		case "servername":
			if server_name, ok := v.(string); ok {
				tls.Enabled = true
				tls.ServerName = server_name
			}
		case "client-fingerprint":
			tls.Enabled = true
			tls.Utls.Enabled = true
			if v, ok := v.(string); ok {
				tls.Utls.Fingerprint = v
			}
		case "tfo":
		case "tls":
		case "ws-headers":
			continue
		case "ws-path":
			continue
		case "encryption":
			continue
		case "cipher":
			continue
		case "alterId":
			continue
		case "udp":
			continue
		case "flow":
			if flow, ok := v.(string); ok {
				outbound["flow"] = flow
			}
		default:
			outbound[k] = v
		}
	}
	if transport.Type != "ws" && transport.Type != "http" && transport.Type != "quic" && transport.Type != "grpc" && transport.Type != "httpupgrade" {
		return outbound
	}
	outbound["transport"] = transport
	outbound["tls"] = tls
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
		transport.Headers = map[string]any{"host": host}
	}
	if transport.Type != "ws" && transport.Type != "http" && transport.Type != "quic" && transport.Type != "grpc" && transport.Type != "httpupgrade" {
		return outbound
	}
	outbound["transport"] = transport
	return outbound
}
