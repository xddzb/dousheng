package model

import (
	"github.com/xddzb/dousheng/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func InitDb() error {
	var err error
	//"root:123456@tcp(127.0.0.1:3306)/dousheng?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := config.DBConnectString()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	//自动建表 会为多对多关系自动建立中间表
	err = db.AutoMigrate(&UserInfo{}, &Video{}, &UserLogin{})
	if err != nil {
		log.Println(err)
	}
	return err
}
