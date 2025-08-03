package nodes

func trojanFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	tls := make(map[string]any)
	tls["enabled"] = true
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
		case "sni":
			tls["server_name"] = v
		case "client-fingerprint":

		default:
			outbound[k] = v
		}
	}
	outbound["tls"] = tls
	return outbound
}
