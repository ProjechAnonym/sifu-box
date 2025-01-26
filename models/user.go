package models

import (
	"github.com/golang-jwt/jwt/v5"
)
type User struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	Email      string `json:"email" yaml:"email"`
	Code       string `json:"code" yaml:"code"`
	PrivateKey string `json:"private_key" yaml:"private_key"`
}

type Jwt struct {
	Admin bool   `json:"admin"`
	UUID  string `json:"uuid"`
	jwt.RegisteredClaims
}

