package singbox

type Log struct {
	Disabled   bool   `json:"disabled" yaml:"disabled"`
	Output     string `json:"output" yaml:"output"`
	Time_stamp bool   `json:"time_stamp" yaml:"time_stamp"`
	Level      string `json:"level" yaml:"level"`
}
