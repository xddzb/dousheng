package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() error {
	var err error
	//username:password@tcp(host:port)/Dbname?charset=utf8&parseTime=True&loc=Local
	dsn := "root:123456@tcp(127.0.0.1:3306)/dousheng?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
