package nodes

import (
	"net/url"
	"strconv"
)

func trojanFromYaml(content map[string]any) map[string]any {
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
		case "password":
			outbound["password"] = v
		case "name":
			outbound["tag"] = v
		case "skip-cert-verify":
			tls.Enabled = true
			tls.Insecure = true
		case "sni":
			tls.Enabled = true
			if server_name, ok := v.(string); ok {
				tls.ServerName = server_name
			}
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
		case "servername":
		case "tfo":
		case "udp":
			continue
		default:
			outbound[k] = v
		}
	}
	outbound["tls"] = tls
	if transport.Type != "ws" && transport.Type != "http" && transport.Type != "quic" && transport.Type != "grpc" && transport.Type != "httpupgrade" {
		return outbound
	}
	outbound["transport"] = transport

	return outbound
}
func trojanFromBase64(content *url.URL) map[string]any {
	outbound := make(map[string]any)
	transport := Transport{}
	tls := TLS{}
	tls.Enabled = true
	tls.ServerName = content.Query().Get("sni")
	if content.Query().Get("allowInsecure") != "1" {
		tls.Insecure = false
	} else {
		tls.Insecure = true
	}
	outbound["tag"] = content.Fragment
	outbound["server"] = content.Hostname()
	port, err := strconv.Atoi(content.Port())
	if err != nil {
		return nil
	}

	outbound["server_port"] = port
	outbound["password"] = content.User.String()
	outbound["type"] = "trojan"
	outbound["tls"] = tls
	transport.Type = content.Query().Get("type")
	if transport.Type == "ws" {
		host := content.Query().Get("host")
		transport.Path = content.Query().Get("path")
		if transport.Type == "ws" {
			transport.Headers = map[string]string{"host": host}
		}
	}
	if transport.Type != "ws" && transport.Type != "http" && transport.Type != "quic" && transport.Type != "grpc" && transport.Type != "httpupgrade" {
		return outbound
	}
	outbound["transport"] = transport

	return outbound
}
