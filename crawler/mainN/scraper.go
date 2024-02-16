package mainN

import (
	"errors"
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/crawler"
	"github.com/SCH-Discord/NoticeNotifier/database"
	"github.com/SCH-Discord/NoticeNotifier/database/model"
	"github.com/SCH-Discord/NoticeNotifier/webhook"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
)

const target = "https://home.sch.ac.kr/sch/06/010100.jsp"

func Scrape() {
	ctx, allocatorCancel, ctxCancel := crawler.CreateCrawler()

	defer allocatorCancel()
	defer ctxCancel()

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(target),
		chromedp.Nodes("#contents_wrap > div > div.board_list > table > tbody > tr", &nodes, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Println(err)
	}

	var db = database.ConnectionDB()
	var latest *model.Latest
	db.Where("notice_type=?", model.MainNotice).Find(&latest)

	isFirst := true
	var title string
	var href string
	var ok bool
	var writer string
	var embeds []webhook.Embed

	for _, node := range nodes {
		err = chromedp.Run(ctx,
			chromedp.Text("a", &title, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("a", "href", &href, &ok, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".writer", &writer, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if err != nil {
			log.Println(err)
			continue
		}
		if !strings.Contains(title, "NEW") {
			continue
		}
		if latest != nil && latest.URL == href {
			break
		}
		if isFirst {
			db.Save(&model.Latest{
				NoticeType: model.MainNotice,
				URL:        href,
			})
			isFirst = true
		}
		embeds = append(embeds, webhook.Embed{
			Title: title,
			Url:   fmt.Sprintf("%s%s", target, href),
			Fields: &[]webhook.Field{
				{
					Name:  "작성자",
					Value: writer,
				},
			},
		})
	}

	if embeds == nil {
		return
	}

	var subscribers []*model.Subscriber
	db.Where("Main = ?", true).Find(&subscribers)

	for _, subscriber := range subscribers {
		go send(subscriber, &embeds)
	}
}

func send(subscriber *model.Subscriber, embeds *[]webhook.Embed) {
	err := webhook.SendMessage(subscriber.URL, &webhook.Message{
		Username: "대학공지",
		Embeds:   embeds,
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
