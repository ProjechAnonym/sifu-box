package models

type Trojan struct {
	Type        string `json:"type"`
	Tag         string `json:"tag"`
	Server      string `json:"server"`
	Server_port int    `json:"server_port"`
	Password    string `json:"password"`
	Tls         *Tls   `json:"tls"`
}