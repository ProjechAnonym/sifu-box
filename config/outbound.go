package config1

type OutboundGroup struct {
	Type      string   `json:"type" yaml:"type"`
	Tag       string   `json:"tag" yaml:"tag"`
	Providers []string `json:"providers" yaml:"providers"`
}
