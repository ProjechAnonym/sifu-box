package model

type Smtp struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Email    string `json:"email" yaml:"email"`
	Password string `json:"password" yaml:"password"`
}
