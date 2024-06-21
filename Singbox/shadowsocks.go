package singbox

import (
	utils "sifu-box/Utils"
)

type shadowsocks struct {
	Type        string `json:"type"`
	Tag         string `json:"tag"`
	Server      string `json:"server"`
	Server_port int    `json:"server_port"`
	Method      string `json:"method"`
	Password    string `json:"password"`
}

func Map_marshal_ss(proxy_map map[string]interface{}) (map[string]interface{}, error) {
	ss := shadowsocks{
		Type:        "ss",
		Tag:         proxy_map["name"].(string),
		Server:      proxy_map["server"].(string),
		Server_port: proxy_map["port"].(int),
		Method:      proxy_map["cipher"].(string),
		Password:    proxy_map["password"].(string),
	}
	ss_map,err := Struct2map(ss,"ss")
	if err != nil {
		utils.Logger_caller("marshal vmess to map failed",err,1)
		return nil,err
	}
	return ss_map,nil
}