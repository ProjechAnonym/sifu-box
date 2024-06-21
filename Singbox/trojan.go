package singbox

import (
	"fmt"
	utils "sifu-box/Utils"
)

type trojan struct {
	Type        string     `json:"type"`
	Tag         string     `json:"tag"`
	Server      string     `json:"server"`
	Server_port int        `json:"server_port"`
	Password    string     `json:"password"`
	Tls         Tls_config `json:"tls"`
}

func Map_marshal_trojan(proxy_map map[string]interface{}) (map[string]interface{}, error) {
	skip_cert_verify,err:=Get_map_value(proxy_map,"skip-cert-verify")
	if err != nil{
		utils.Logger_caller(fmt.Sprintf("%s has no skip-cert-verify key",proxy_map["name"].(string)),err,1)
		skip_cert_verify = false
	}
	sni,err := Get_map_value(proxy_map,"sni")
	if err != nil{
		utils.Logger_caller(fmt.Sprintf("%s has no sni key",proxy_map["sni"].(string)),err,1)
		return nil, err
	}
	trojan := trojan{
		Type:        "trojan",
		Tag:         proxy_map["name"].(string),
		Server:      proxy_map["server"].(string),
		Server_port: proxy_map["port"].(int),
		Password:    proxy_map["password"].(string),
		Tls: Tls_config{
			Enabled: true,
			Insecure: skip_cert_verify.(bool),
			Server_name: sni.(string),
		},
	}
	trojan_map,err := Struct2map(trojan,"trojan")
	if err != nil{
		utils.Logger_caller("marshal trojan to map failed",err,1)
		return nil, err
	}
	return trojan_map, nil
}
