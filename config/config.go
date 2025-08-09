package config1

type Config struct {
	Experiment *Experiment      `json:"experiment,omitempty" yaml:"experiment,omitempty"`
	Ntp        *Ntp             `json:"ntp,omitempty" yaml:"ntp,omitempty"`
	Log        *Log             `json:"log,omitempty" yaml:"log,omitempty"`
	DNS        DNS              `json:"dns" yaml:"dns"`
	Inbounds   []Inbound        `json:"inbounds" yaml:"inbounds"`
	Outbounds  []map[string]any `json:"outbounds" yaml:"outbounds"`
	Route      Route            `json:"route" yaml:"route"`
}
