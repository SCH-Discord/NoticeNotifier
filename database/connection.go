package database

import (
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var connDB *gorm.DB = nil

func ConnectionDB() *gorm.DB {
	if connDB != nil {
		db, err := connDB.DB()
		if err == nil && db.Ping() == nil {
			return connDB
		}
	}
	dsn := fmt.Sprintf("root:%s@tcp(mariadb)/%s?parseTime=True", config.DBPassword(), config.DBName())
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		connDB = db
		return db
	}
	log.Fatalln("err")
	return nil
}
