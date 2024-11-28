package models

type Host struct {
	ID          uint64 `gorm:"primaryKey" json:"-" yaml:"-"`
	Url         string `json:"url" yaml:"url" gorm:"unique;not null;type:varchar(255)"`
	Username    string `json:"username" yaml:"username" gorm:"not null;type:varchar(20)"`
	Password    string `json:"password" yaml:"password" gorm:"not null;type:varchar(255)"`
	Localhost   bool   `json:"localhost" yaml:"localhost" gorm:"not null;type:bool"`
	Config      string `json:"config" yaml:"config" gorm:"type:varchar(255)"`
	Fingerprint string `json:"-" yaml:"fingerprint" gorm:"type:varchar(255)"`
	Secret      string `json:"secret" yaml:"secret" gorm:"type:varchar(255)"`
	Port        uint16 `json:"port" yaml:"port" gorm:"not null;type:uint16"`
	Template    string `json:"template" yaml:"template" gorm:"not null;type:varchar(255)"`
}