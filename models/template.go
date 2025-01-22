package models

import "go.uber.org/zap"

type Outbound interface {
	Transform(message map[string]interface{}, logger *zap.Logger) (Outbound, error)
	GetTag() string
}

type Template struct {
	Log          *Log             `json:"log,omitempty" yaml:"log,omitempty"`
	Ntp          *NTP             `json:"ntp,omitempty" yaml:"ntp,omitempty"`
	Experimental *Experimental    `json:"experimental,omitempty" yaml:"experimental,omitempty"`
	Inbounds     []Inbounds       `json:"inbounds" yaml:"inbounds"`
	Dns          DNS              `json:"dns" yaml:"dns"`
	Route        Route            `json:"route" yaml:"route"`
	Outbounds    []interface{} 	  `json:"outbounds" yaml:"outbounds"`
	CustomOutbounds []interface{} `json:"custom_outbounds,omitempty" yaml:"custom_outbounds,omitempty"`
}

func (t *Template) SetOutbounds(outbounds []Outbound) {
	interfaceOutbounds := make([]interface{}, len(outbounds) + len(t.Outbounds) + len(t.CustomOutbounds))
	t.Outbounds = append(t.Outbounds, t.CustomOutbounds...)
	copy(interfaceOutbounds, t.Outbounds)
	for i, outbound := range outbounds {
		interfaceOutbounds[i + len(t.Outbounds)] = outbound
	}
	t.Outbounds = interfaceOutbounds
	t.CustomOutbounds = nil
}
