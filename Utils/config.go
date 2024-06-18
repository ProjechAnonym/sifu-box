package utils

type Box_url struct {
	Url   string `yaml:"url"`
	Proxy bool   `yaml:"proxy"`
	Label string `yaml:"label"`
}
type Ruleset_value struct {
	Path            string `yaml:"path"`
	Format          string `yaml:"format"`
	Type            string `yaml:"type"`
	China           bool   `yaml:"china"`
	Update_interval string `yaml:"update_interval"`
	Download_detour string `yaml:"download_detour"`
}
type Box_ruleset struct {
	Label string        `yaml:"label"`
	Value Ruleset_value `yaml:"value"`
}

type Box_config struct {
	Url      []Box_url     `yaml:"url"`
	Rule_set []Box_ruleset `yaml:"rule_set"`
}
type Cors struct {
	Origins []string `yaml:"origins"`
}
type Server_config struct {
	Cors        Cors   `yaml:"cors"`
	Key         string `yaml:"key"`
	Server_mode bool   `yaml:"server_mode"`
}