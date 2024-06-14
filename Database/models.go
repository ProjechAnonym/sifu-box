package database

type Server struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	Url       string `json:"url" gorm:"unique;not null;type:varchar(255)"`
	Username  string `json:"username" gorm:"not null;type:varchar(20)"`
	Password  string `json:"password" gorm:"not null;type:varchar(255)"`
	Localhost bool   `json:"localhost" gorm:"not null;type:bool"`
}