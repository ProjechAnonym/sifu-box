package nodes

import (
	"fmt"
	"net/url"
	"strconv"
)

func hysteria2FromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	tls := TLS{}
	obfs := make(map[string]any)
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
				if server_name == "" {
					tls.Enabled = false
					continue
				}
				tls.ServerName = server_name
			} else {
				tls.Enabled = false
				continue
			}
		case "obfs":
			obfs["type"] = v
		case "obfs-password":
			obfs["password"] = v
		case "udp":
			continue
		case "down":
			continue
		case "up":
			continue
		default:
			outbound[k] = v
		}
	}
	outbound["obfs"] = obfs
	if tls.Enabled {
		outbound["tls"] = tls
	}
	return outbound
}

func hysteria2FromBase64(content *url.URL) map[string]any {
	fmt.Println(content)
	outbound := make(map[string]any)
	obfs := make(map[string]any)
	tls := TLS{}
	outbound["tag"] = content.Fragment
	outbound["server"] = content.Hostname()
	outbound["password"] = content.User.Username()
	port, err := strconv.Atoi(content.Port())
	if err != nil {
		return nil
	}
	outbound["server_port"] = port
	outbound["type"] = "hysteria2"
	if content.Query().Get("insecure") == "0" {
		tls.Enabled = false
	}
	tls.ServerName = content.Query().Get("sni")
	tls.Alpn = []string{content.Query().Get("alpn")}
	obfs["type"] = content.Query().Get("obfs")
	obfs["password"] = content.Query().Get("obfs-password")
	outbound["obfs"] = obfs
	if tls.Enabled {
		outbound["tls"] = tls
	}
	return outbound
}
