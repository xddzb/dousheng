package model

import (
	"github.com/xddzb/dousheng/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() error {
	var err error
	//"root:123456@tcp(127.0.0.1:3306)/dousheng?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := config.DBConnectString()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
