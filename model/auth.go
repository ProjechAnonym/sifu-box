package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Code     string `json:"code" yaml:"code"`
	Key      string `json:"key" yaml:"key"`
}

type Jwt struct {
	Admin bool   `json:"admin"`
	UUID  string `json:"uuid"`
	jwt.RegisteredClaims
}
