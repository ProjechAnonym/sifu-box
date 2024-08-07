package models

type ShadowSocks struct {
	Type        string `json:"type"`
	Tag         string `json:"tag"`
	Server      string `json:"server"`
	Server_port int    `json:"server_port"`
	Method      string `json:"method"`
	Password    string `json:"password"`
}