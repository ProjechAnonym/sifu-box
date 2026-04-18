package nodes

import (
	"net/url"
	"strconv"
)

func tuicFromYaml(content map[string]any) map[string]any {
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
		case "uuid":
			outbound["uuid"] = v
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
		case "udp":
			continue
		case "reduce-rtt":
			continue
		default:
			outbound[k] = v
		}
	}
	outbound["tls"] = tls
	return outbound
}

func tuicFromBase64(content *url.URL) map[string]any {
	outbound := make(map[string]any)
	tls := TLS{}
	outbound["tag"] = content.Fragment
	outbound["server"] = content.Hostname()
	outbound["uuid"] = content.User.Username()
	password, tag := content.User.Password()
	if !tag {
		return nil
	}
	outbound["password"] = password
	port, err := strconv.Atoi(content.Port())
	if err != nil {
		return nil
	}
	outbound["server_port"] = port
	outbound["type"] = "tuic"
	if content.Query().Get("allow_insecure") == "1" {
		tls.Enabled = true
		tls.Insecure = true
	}
	tls.ServerName = content.Query().Get("sni")
	tls.Alpn = []string{content.Query().Get("alpn")}
	outbound["congestion_control"] = "bbr"
	outbound["udp_relay_mode"] = content.Query().Get("udp_relay_mode")
	outbound["tls"] = tls
	return outbound
}
