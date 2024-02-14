package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Ctrl + C를 통해 종료할 수 있습니다.")

	for {
		select {
		case <-done:
			fmt.Println("프로그램을 종료합니다.")
			os.Exit(0)
		case <-time.After(timeUntilNextRun()):
			doTask()
		}
	}
}

func doTask() {
	//TODO
	fmt.Println("test")
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
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), nextRunHour(now.Hour()), 0, 0, 0, now.Location())

	if now.After(targetTime) {
		targetTime = targetTime.Add(24 * time.Hour)
	}

	log.Println("다음실행 시간:", targetTime)

	return targetTime.Sub(now)
}
