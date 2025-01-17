package models

type NTP struct {
	Enabled    bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Server     string `json:"server,omitempty" yaml:"server,omitempty"`
	ServerPort int    `json:"server_port,omitempty" yaml:"server_port,omitempty"`
	Interval   string `json:"interval,omitempty" yaml:"interval,omitempty"`
	Dial
}