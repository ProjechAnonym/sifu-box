package config1

type Ntp struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Server      string `json:"server" yaml:"server"`
	Server_port int    `json:"server_port,omitempty" yaml:"server_port,omitempty"`
	Interval    string `json:"interval,omitempty" yaml:"interval,omitempty"`
}
