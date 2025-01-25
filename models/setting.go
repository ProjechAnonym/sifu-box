package models

type Server struct {
	Interval string `json:"interval,omitempty" yaml:"interval,omitempty"`
	SSL      *struct { 
		Public  string `json:"public" yaml:"public"`
		Private string `json:"private" yaml:"private"`
	} `json:"ssl,omitempty" yaml:"ssl,omitempty"`
}
type Command struct {
	Name  string        `json:"name" yaml:"name"`
	Args  []interface{} `json:"args,omitempty" yaml:"args,omitempty"`
}
type Singbox struct {
	Template   string `json:"template" yaml:"template"`
	Provider   string `json:"provider" yaml:"provider"`
	WorkDir    string `json:"work_dir" yaml:"work_dir"`
	ConfigPath string `json:"config_path" yaml:"config_path"`
	BinaryPath string `json:"binary_path" yaml:"binary_path"`
	Commands    map[string]*Command `json:"commands" yaml:"commands"`
}
type Application struct {
	Server  *Server  `json:"server,omitempty" yaml:"server,omitempty"`
	Singbox *Singbox `json:"singbox,omitempty" yaml:"singbox,omitempty"`
}
type Setting struct {
	Application *Application `json:"application,omitempty" yaml:"application,omitempty"`
	Configuration *Configuration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
}

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
	NameServer     string `json:"name_server,omitempty" yaml:"name_server,omitempty"`
	Label          string `json:"label" yaml:"label"`
	DownloadDetour string `json:"download_detour" yaml:"download_detour"`
	UpdateInterval string `json:"update_interval" yaml:"update_interval"`
}

type Configuration struct {
	Providers []Provider          `json:"providers,omitempty" yaml:"providers,omitempty"`
	Rulesets  []RuleSet           `json:"rulesets,omitempty" yaml:"rulesets,omitempty"`
	Templates map[string]Template `json:"templates,omitempty" yaml:"templates,omitempty"`
}