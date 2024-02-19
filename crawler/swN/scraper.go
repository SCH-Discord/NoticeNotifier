package swN

import (
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/crawler"
	"github.com/SCH-Discord/NoticeNotifier/database"
	"github.com/SCH-Discord/NoticeNotifier/database/model"
	"github.com/SCH-Discord/NoticeNotifier/webhook"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"math/big"
	"time"
)

const target = "https://home.sch.ac.kr/sw/07/010000.jsp"

func scrape(name string) {
	ctx, allocatorCancel, ctxCancel, timeoutCancel := crawler.CreateCrawler()

	defer allocatorCancel()
	defer ctxCancel()
	defer timeoutCancel()

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(target),
		chromedp.WaitVisible("#sub_board > div > div.board_list > table > tbody", chromedp.ByQuery),
		chromedp.Nodes("#sub_board > div > div.board_list > table > tbody > tr", &nodes, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Println(err)
		return
	}

	var db = database.ConnectionDB()
	var mLatest *model.Latest
	db.Where("notice_type=?", model.SWNotice).Find(&mLatest)

	var latest big.Int
	if mLatest != nil {
		latest.SetString(mLatest.URL, 10)
	}

	nowDate := crawler.NowDate()

	var newLatest big.Int
	var articleNo big.Int
	var articleNoStr string
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
		_, err := fmt.Sscanf(href, "?mode=view&article_no=%s", &articleNoStr)
		if err != nil {
			log.Println(err)
			continue
		}
		articleNo.SetString(articleNoStr, 10)
		if cmp := latest.Cmp(&articleNo); cmp == 0 || cmp == 1 {
			continue
		} else if newLatest.Cmp(&articleNo) == -1 {
			newLatest.SetString(articleNoStr, 10)
		}

		embeds = append(embeds, webhook.Embed{
			Title: crawler.FixTitle(title),
			Url:   fmt.Sprintf("%s%s", target, href),
			Fields: &[]webhook.Field{
				{
					Name:  "작성자",
					Value: writer,
				},
			},
		})
	}

	if newLatest.String() != "0" {
		db.Save(&model.Latest{
			NoticeType: model.SWNotice,
			URL:        newLatest.String(),
		})
	} else {
		return
	}

	var subscribers []*model.Subscriber
	db.Where("SW = ?", true).Find(&subscribers)

	for _, subscriber := range subscribers {
		go crawler.Send(name, subscriber, &embeds)
	}
}

func Scrape() {
	log.Println("Start scrape SW중심대학산업단")
	scrape("SW중심대학산업단")
	time.Sleep(crawler.WaitTime)
	log.Println("Completed scrape SW중심대학산업단")
}
