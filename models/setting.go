package models

type Provider struct {
	Name   string `json:"name" yaml:"name"`
	Path   string `json:"path" yaml:"path"`
	Remote bool   `json:"remote" yaml:"remote"`
	Detour string `json:"detour" yaml:"detour"`
}
type RuleSet struct {
	Tag            string `json:"tag" yaml:"tag"`
	Type           string `json:"type" yaml:"type"`
	Path           string `json:"path" yaml:"path"`
	Format         string `json:"format" yaml:"format"`
	China          bool   `json:"china" yaml:"china"`
	Outbound       string `json:"outbound" yaml:"outbound"`
	Label          string `json:"label" yaml:"label"`
	DownloadDetour string `json:"download_detour" yaml:"download_detour"`
	UpdateInterval string `json:"update_interval" yaml:"update_interval"`
}
type Server struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}
type Singbox struct {
	WorkDir    string `json:"work_dir" yaml:"work_dir"`
	ConfigPath string `json:"config_path" yaml:"config_path"`
	BinaryPath string `json:"binary_path" yaml:"binary_path"`
}
type Setting struct {
	Providers []Provider          `json:"providers,omitempty" yaml:"providers,omitempty"`
	Rulesets  []RuleSet           `json:"rulesets,omitempty" yaml:"rulesets,omitempty"`
	Templates map[string]Template `json:"templates,omitempty" yaml:"templates,omitempty"`
	Server    Server              `json:"server" yaml:"server"`
	Singbox   *Singbox            `json:"singbox" yaml:"singbox"`
}