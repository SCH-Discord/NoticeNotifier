package crawler

import (
	"errors"
	"github.com/SCH-Discord/NoticeNotifier/database"
	"github.com/SCH-Discord/NoticeNotifier/database/model"
	"github.com/SCH-Discord/NoticeNotifier/webhook"
	"log"
)

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
