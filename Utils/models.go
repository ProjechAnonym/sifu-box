package utils

type Server struct {
	ID          uint64 `gorm:"primaryKey" json:"-"`
	Url         string `json:"url" gorm:"unique;not null;type:varchar(255)"`
	Username    string `json:"username" gorm:"not null;type:varchar(20)"`
	Password    string `json:"password" gorm:"not null;type:varchar(255)"`
	Localhost   bool   `json:"localhost" gorm:"not null;type:bool"`
	Config      string `json:"config" gorm:"type:varchar(255)"`
	Fingerprint string `json:"-" gorm:"type:varchar(255)"`
	Secret      string `json:"secret" gorm:"type:varchar(255)"`
	Port        uint16 `json:"port" gorm:"not null;type:uint16"`
}