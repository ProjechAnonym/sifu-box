package nodes

import (
	"net/url"
	"strconv"
)

func anytlsFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
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
		case "server":
			outbound["server"] = v
		case "skip-cert-verify":
			tls.Enabled = true
			tls.Insecure = true
		case "sni":
			tls.Enabled = true
			if server_name, ok := v.(string); ok {
				tls.ServerName = server_name
			}
		case "alpn":
			if alpn, ok := v.([]string); ok {
				tls.Alpn = alpn
			}
		case "client-fingerprint":
			continue
		case "udp":
			continue
		default:
			outbound[k] = v
		}
	}
	outbound["tls"] = tls
	return outbound
}

func anytlsFromBase64(content *url.URL) map[string]any {
	tls := TLS{}
	outbound := make(map[string]any)
	outbound["tag"] = content.Fragment
	outbound["server"] = content.Hostname()
	port, err := strconv.Atoi(content.Port())
	if err != nil {
		return nil
	}
	outbound["server_port"] = port
	outbound["password"] = content.User.String()
	if content.Query().Get("insecure") == "1" {
		tls.Enabled = true
		tls.Insecure = true
		tls.ServerName = content.Query().Get("sni")
	} else {
		tls.Enabled = true
		tls.ServerName = content.Query().Get("sni")
	}
	outbound["type"] = "anytls"
	outbound["tls"] = tls
	return outbound
}
