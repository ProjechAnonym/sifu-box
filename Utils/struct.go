package utils

type Box_url struct {
	Path   string `yaml:"path" json:"path"`
	Proxy  bool   `yaml:"proxy" json:"proxy"`
	Label  string `yaml:"label" json:"label"`
	Remote bool   `yaml:"remote" json:"remote"`
}
type Ruleset_value struct {
	Path            string `yaml:"path" json:"path"`
	Format          string `yaml:"format" json:"format"`
	Type            string `yaml:"type" json:"type"`
	China           bool   `yaml:"china" json:"china"`
	Update_interval string `yaml:"update_interval" json:"update_interval"`
	Download_detour string `yaml:"download_detour" json:"download_detour"`
}
type Box_ruleset struct {
	Label string        `yaml:"label" json:"label"`
	Value Ruleset_value `yaml:"value" json:"value"`
}

type Box_config struct {
	Url      []Box_url     `yaml:"url" json:"url"`
	Rule_set []Box_ruleset `yaml:"rule_set" json:"rule_set"`
}
type Cors struct {
	Origins []string `yaml:"origins"`
}
type Server_config struct {
	Cors        Cors   `yaml:"cors"`
	Key         string `yaml:"key"`
	Server_mode bool   `yaml:"server_mode"`
	Token 		string `yaml:"token"`
}