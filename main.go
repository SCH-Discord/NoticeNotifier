package main

import (
	"database/sql"
	"github.com/SCH-Discord/NoticeNotifier/crawler/libraryN"
	"github.com/SCH-Discord/NoticeNotifier/crawler/mainN"
	"github.com/SCH-Discord/NoticeNotifier/database"
	"github.com/SCH-Discord/NoticeNotifier/database/model"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	sqlDb, err := setupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDb.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Ctrl + C를 통해 종료할 수 있습니다.")

loop:
	for {
		select {
		case <-done:
			log.Println("프로그램을 종료합니다.")
			break loop
		case <-time.After(timeUntilNextRun()):
			go doTask()
		}
	}
}

func doTask() {
	log.Println("Starting a scheduled task")
	mainN.Scrape()
	libraryN.Scrape()
}

// 데이터 베이스 설정
func setupDatabase() (*sql.DB, error) {
	db := database.ConnectionDB()
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Subscriber{})
	if err != nil {
		sqlDb.Close()
		return nil, err
	}

	err = db.AutoMigrate(&model.Latest{})
	if err != nil {
		sqlDb.Close()
		return nil, err
	}

	log.Println("ORM ready")
	return sqlDb, nil
}

// 다음 실행 시간 구하기
var targetHours = []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

func nextRunHour(currentHour int) int {
	for _, h := range targetHours {
		if h > currentHour {
			return h
		}
	}
	return targetHours[0]
}

// 다음 실행 시간 까지 남은 시간리턴
func timeUntilNextRun() time.Duration {
	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), nextRunHour(now.Hour()), rand.Intn(5), rand.Intn(59), now.Nanosecond(), now.Location())

	if now.After(targetTime) {
		targetTime = targetTime.Add(24 * time.Hour)
	}

	log.Println("다음실행 시간:", targetTime)

	return targetTime.Sub(now)
}
