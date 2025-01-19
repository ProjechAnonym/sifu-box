package models

type Outbounds interface{}

type Template struct {
	Log          *Log          `json:"log,omitempty" yaml:"log,omitempty"`
	Ntp          *NTP          `json:"ntp,omitempty" yaml:"ntp,omitempty"`
	Experimental *Experimental `json:"experimental,omitempty" yaml:"experimental,omitempty"`
	Inbounds     []Inbounds    `json:"inbounds" yaml:"inbounds"`
	Dns          DNS           `json:"dns" yaml:"dns"`
	Route        Route         `json:"route" yaml:"route"`
	Outbounds    []Outbounds   `json:"outbounds" yaml:"outbounds"`
	UDP          bool          `json:"udp,omitempty" yaml:"udp,omitempty"`
}
