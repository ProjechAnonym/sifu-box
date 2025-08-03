package nodes

func vmessFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	transport := make(map[string]any)
	for k, v := range content {
		switch k {
		case "port":
			outbound["server_port"] = v
		case "cipher":
			outbound["security"] = v
		case "name":
			outbound["tag"] = v
		case "alterId":
			outbound["alter_id"] = v
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
		default:
			outbound[k] = v
		}
	}
	outbound["transport"] = transport
	return outbound

}
