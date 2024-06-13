package database

import (
	"fmt"
	utils "sifu-box/Utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)
var Db *gorm.DB
func Get_database(){
	project_dir,_ := utils.Get_value("project-dir")
	Db,_ = gorm.Open(sqlite.Open(fmt.Sprintf("%s/sifu-box.db",project_dir)),&gorm.Config{})
	Db.AutoMigrate(&Server{})
}