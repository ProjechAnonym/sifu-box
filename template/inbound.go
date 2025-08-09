package template

type Inbound struct {
	Type  string         `json:"type" yaml:"type"`
	Tag   string         `json:"tag" yaml:"tag"`
	Extra map[string]any `json:",inline" yaml:",inline"`
}
