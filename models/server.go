package models

type Server struct {
	Mode   bool     `yaml:"mode"`
	Token  string   `yaml:"token"`
	Cors   []string `yaml:"cors"`
	Key    string   `yaml:"key"`
	Listen string   `yaml:"listen"`
}