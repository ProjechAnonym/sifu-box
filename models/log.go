package models

type Log struct {
	Level     string `json:"level,omitempty" yaml:"level,omitempty"`
	Output    string `json:"output,omitempty" yaml:"output,omitempty"`
	Disabled  bool   `json:"disabled" yaml:"disabled"`
	TimeStamp bool   `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
}