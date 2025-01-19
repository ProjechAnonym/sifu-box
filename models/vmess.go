package models

type VMess struct {
	Transport           *Transport             `json:"transport,omitempty" yaml:"transport,omitempty"`
	Multiplex           *Multiplex             `json:"multiplex,omitempty" yaml:"multiplex,omitempty"`
	UDP                 bool                   `json:"-" yaml:"udp,omitempty"`
	TLS                 map[string]interface{} `json:"tls,omitempty" yaml:"tls,omitempty"`
	Type                string                 `json:"type" yaml:"type"`
	Tag                 string                 `json:"tag" yaml:"name"`
	Server              string                 `json:"server" yaml:"server"`
	ServerPort          int                    `json:"server_port" yaml:"port"`
	Network             string                 `json:"network,omitempty" yaml:"network,omitempty"`
	UUID                string                 `json:"uuid" yaml:"uuid"`
	Security            string                 `json:"security,omitempty" yaml:"cipher,omitempty"`
	AlterID             int                    `json:"alter_id" yaml:"alterId"`
	GlobalPadding       bool                   `json:"global_padding,omitempty" yaml:"global_padding,omitempty"`
	AuthenticatedLength bool                   `json:"authenticated_length,omitempty" yaml:"authenticated_length,omitempty"`
	PacketEncoding      string                 `json:"packet_encoding,omitempty" yaml:"packet_encoding,omitempty"`
}
