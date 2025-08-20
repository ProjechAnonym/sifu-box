package singbox

type Log struct {
	Disabled   bool   `json:"disabled" yaml:"disabled"`
	Output     string `json:"output" yaml:"output"`
	Time_stamp bool   `json:"timestamp" yaml:"timestamp"`
	Level      string `json:"level" yaml:"level"`
}
