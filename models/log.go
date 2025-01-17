package models

type Log struct {
	Level     string `json:"level,omitempty" yaml:"level,omitempty"`
	Output    string `json:"output,omitempty" yaml:"output,omitempty"`
	Disabled  bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	TimeStamp bool   `json:"time_stamp,omitempty" yaml:"time_stamp,omitempty"`
}