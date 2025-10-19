package singbox

type Log struct {
	Disabled   bool   `json:"disabled" yaml:"disabled"`
	Output     string `json:"output,omitempty" yaml:"output,omitempty"`
	Time_stamp bool   `json:"timestamp" yaml:"timestamp"`
	Level      string `json:"level,omitempty" yaml:"level,omitempty"`
}
