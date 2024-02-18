package crawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/database"
	"github.com/SCH-Discord/NoticeNotifier/database/model"
	"github.com/SCH-Discord/NoticeNotifier/webhook"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

const WaitTime = 50 * time.Millisecond

func CreateCrawler() (context.Context, context.CancelFunc, context.CancelFunc) {
	allocatorContext, allocatorCancel := chromedp.NewRemoteAllocator(context.Background(), "wss://chrome:9222")
	ctx, ctxCancel := chromedp.NewContext(allocatorContext)

	return ctx, allocatorCancel, ctxCancel
}

func Send(name string, subscriber *model.Subscriber, embeds *[]webhook.Embed) {
	err := webhook.SendMessage(subscriber.URL, &webhook.Message{
		Username:  name,
		AvatarUrl: "https://raw.githubusercontent.com/SCH-Discord/image/main/profile.png",
		Embeds:    embeds,
	})
	if err == nil {
		return
	}
	var notOk *webhook.NotOk
	if errors.As(err, &notOk) {
		log.Printf("remove %s\n", subscriber.URL)
		database.ConnectionDB().Delete(subscriber)
	} else {
		log.Println(err)
	}
}

var day int
var date string

func NowDate() string {
	now := time.Now()
	if date == "" {
		day = now.Day()
		date = fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	} else if day != now.Day() {
		day = now.Day()
		date = fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	}
	return date
}
