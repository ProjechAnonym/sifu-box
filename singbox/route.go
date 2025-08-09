package singbox

type Route struct {
	Rules                    []map[string]any `json:"rules,omitempty" yaml:"rules,omitempty"`
	Rule_sets                []Rule_set       `json:"rule_set,omitempty" yaml:"rule_set,omitempty"`
	Final                    string           `json:"final,omitempty" yaml:"final,omitempty"`
	Auto_detect_interface    bool             `json:"auto_detect_interface,omitempty" yaml:"auto_detect_interface,omitempty"`
	Override_android_vpn     bool             `json:"override_android_vpn,omitempty" yaml:"override_android_vpn,omitempty"`
	Default_interface        string           `json:"default_interface,omitempty" yaml:"default_interface,omitempty"`
	Default_mark             int              `json:"default_mark,omitempty" yaml:"default_mark,omitempty"`
	Default_network_strategy string           `json:"default_network_strategy,omitempty" yaml:"default_network_strategy,omitempty"`
	Default_fallback_delay   string           `json:"default_fallback_delay,omitempty" yaml:"default_fallback_delay,omitempty"`
}
type Rule_set struct {
	Type            string `json:"type" yaml:"type"`
	Tag             string `json:"tag" yaml:"tag"`
	Format          string `json:"format" yaml:"format"`
	URL             string `json:"url,omitempty" yaml:"url,omitempty"`
	Download_detour string `json:"download_detour,omitempty" yaml:"download_detour,omitempty"`
	Update_interval string `json:"update_interval,omitempty" yaml:"update_interval,omitempty"`
	Path            string `json:"path,omitempty" yaml:"path,omitempty"`
}
