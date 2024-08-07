package models

type Vmess struct {
	Tag         string      `json:"tag"`
	Server      string      `json:"server"`
	Server_port int         `json:"server_port"`
	Alter_id    int         `json:"alter_id"`
	Type        string      `json:"type"`
	Uuid        string      `json:"uuid"`
	Security    string      `json:"security"`
	Network     string      `json:"network,omitempty"`
	Tls         *Tls        `json:"tls,omitempty"`
	Transport   interface{} `json:"transport,omitempty"`
}