package nodes

func shadowsocksFromYaml(content map[string]any) map[string]any {
	outbound := make(map[string]any)
	for k, v := range content {
		switch k {
		case "port":
			outbound["server_port"] = v
		case "cipher":
			outbound["method"] = v
		case "name":
			outbound["tag"] = v
		case "type":
			outbound["type"] = "shadowsocks"
		case "udp":
		default:
			outbound[k] = v
		}
	}
	return outbound
}
