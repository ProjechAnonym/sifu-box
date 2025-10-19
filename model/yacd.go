package model

type Yacd struct {
	Url      string `json:"url" yaml:"url"`
	Secret   string `json:"secret" yaml:"secret"`
	Template string `json:"template,omitempty" yaml:"template,omitempty"`
	Log      bool   `json:"log" yaml:"log"`
}
