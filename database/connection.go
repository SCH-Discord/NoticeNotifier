package database

import (
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var connDB *gorm.DB = nil

func ConnectionDB() (*gorm.DB, error) {
	if connDB != nil {
		db, err := connDB.DB()
		if err == nil && db.Ping() == nil {
			return connDB, nil
		}
	}
	dsn := fmt.Sprintf("root:%s@tcp(mariadb)/%s?parseTime=True", config.DBPassword(), config.DBName())
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		connDB = db
		return db, nil
	}
	log.Fatalln("DB 연결에 실패했습니다.")
	return nil, err
}
