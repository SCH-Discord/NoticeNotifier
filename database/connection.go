package database

import (
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var connDB *gorm.DB = nil

// DB 연결
func ConnectionDB() *gorm.DB {
	if connDB != nil {
		db, err := connDB.DB()
		if err == nil && db.Ping() == nil {
			return connDB
		}
	}

	const maxRetries = 10

	for i := 0; i < maxRetries; i++ {
		dsn := fmt.Sprintf("root:%s@tcp(mariadb)/%s?parseTime=True", config.DBPassword(), config.DBName())
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			connDB = db
			return db
		}
		time.Sleep(100 * time.Millisecond)
	}

	log.Fatalln("Failed to connect to the database after", maxRetries, "attempts")
	return nil
}
