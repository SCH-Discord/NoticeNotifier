package mainN

import (
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/crawler"
	"github.com/SCH-Discord/NoticeNotifier/database"
	"github.com/SCH-Discord/NoticeNotifier/database/model"
	"github.com/SCH-Discord/NoticeNotifier/webhook"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
)

const target = "https://home.sch.ac.kr/sch/06/%s.jsp"

const university = "010100"  // 대학공지
const academic = "010200"    // 학사공지
const scholarship = "010300" // 장학공지

func scrape(code string, name string) {
	ctx, allocatorCancel, ctxCancel := crawler.CreateCrawler()

	defer allocatorCancel()
	defer ctxCancel()

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf(target, code)),
		chromedp.Nodes("#contents_wrap > div > div.board_list > table > tbody > tr", &nodes, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Println(err)
	}

	var db = database.ConnectionDB()
	var latest *model.Latest
	db.Where("notice_type=?", fmt.Sprintf("%s%s", model.MainNotice, code)).Find(&latest)

	nowDate := crawler.NowDate()

	isFirst := true
	var title string
	var href string
	var ok bool
	var writer string
	var date string
	var embeds []webhook.Embed

	for _, node := range nodes {
		err = chromedp.Run(ctx,
			chromedp.Text("a", &title, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("a", "href", &href, &ok, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".writer", &writer, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".date", &date, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if err != nil {
			log.Println(err)
			continue
		}
		if date != nowDate {
			continue
		}
		if latest != nil && latest.URL == href {
			break
		}
		if isFirst {
			db.Save(&model.Latest{
				NoticeType: fmt.Sprintf("%s%s", model.MainNotice, code),
				URL:        href,
			})
			isFirst = false
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
		go crawler.Send(name, subscriber, &embeds)
	}
}

func Scrape() {
	scrape(university, "대학공지")
	scrape(academic, "학사공지")
	scrape(scholarship, "장학공지")
}
