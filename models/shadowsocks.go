package models

type ShadowSocks struct {
	Type       string `json:"type" yaml:"type"`
	Tag        string `json:"tag" yaml:"name"`
	Server     string `json:"server" yaml:"server"`
	ServerPort int    `json:"server_port" yaml:"port"`
	UDP        bool   `json:"-" yaml:"udp,omitempty"`
	Method     string `json:"method" yaml:"cipher"`
	Password   string `json:"password" yaml:"password"`
	Plugin     string `json:"plugin,omitempty" yaml:"plugin,omitempty"`
	PluginOpts string `json:"plugin_opts,omitempty" yaml:"plugin_opts,omitempty"`
	Dial       `json:",inline" yaml:",inline"`
}