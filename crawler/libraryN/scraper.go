package libraryN

import (
	"fmt"
	"github.com/SCH-Discord/NoticeNotifier/crawler"
	"github.com/SCH-Discord/NoticeNotifier/database"
	"github.com/SCH-Discord/NoticeNotifier/database/model"
	"github.com/SCH-Discord/NoticeNotifier/webhook"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

const root = "https://library.sch.ac.kr"
const target = "https://library.sch.ac.kr/bbs/list/%d"

const general = 1  // 일반공지
const academic = 2 // 학술공지
const event = 3    // 교육/행사공지

func scrape(code int, name string) {
	ctx, allocatorCancel, ctxCancel := crawler.CreateCrawler()

	defer allocatorCancel()
	defer ctxCancel()

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf(target, code)),
		chromedp.WaitReady("#divContent > form > div > table > tbody", chromedp.ByQuery),
		chromedp.Nodes("#divContent > form > div > table > tbody > tr", &nodes, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Println(err)
	}

	var db = database.ConnectionDB()
	var mLatest *model.Latest
	db.Where("notice_type=?", fmt.Sprintf("%s%d", model.LibraryNotice, code)).Find(&mLatest)

	var latest int64 = 0
	if mLatest != nil {
		_, err = fmt.Sscanf(mLatest.URL, "%d", &latest)
		if err != nil {
			latest = 0
		}
	}

	nowDate := crawler.NowDate()

	var newLatest int64 = 0
	var postNo int64
	var title string
	var href string
	var ok bool
	var date string
	var embeds []webhook.Embed

	for _, node := range nodes {
		err = chromedp.Run(ctx,
			chromedp.Text("a", &title, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("a", "href", &href, &ok, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(".reportDate", &date, chromedp.ByQuery, chromedp.FromNode(node)),
		)

		if err != nil {
			log.Println(err)
			continue
		}
		if date != nowDate {
			continue
		}
		_, err := fmt.Sscanf(href, fmt.Sprintf("/bbs/content/%d_%%d", code), &postNo)
		if err != nil {
			log.Println(err)
			continue
		}
		if latest >= postNo {
			continue
		} else if newLatest < postNo {
			newLatest = postNo
		}

		embeds = append(embeds, webhook.Embed{
			Title: title,
			Url:   fmt.Sprintf("%s%s", root, href),
		})
	}

	if newLatest != 0 {
		db.Save(&model.Latest{
			NoticeType: fmt.Sprintf("%s%d", model.LibraryNotice, code),
			URL:        fmt.Sprintf("%d", newLatest),
		})
	} else {
		return
	}

	var subscribers []*model.Subscriber
	db.Where("Library = ?", true).Find(&subscribers)

	for _, subscriber := range subscribers {
		go crawler.Send(name, subscriber, &embeds)
	}
}

func Scrape() {
	log.Println("Start scrape 도서관")
	log.Println("도서관(1/3)")
	scrape(general, "도서관")
	time.Sleep(crawler.WaitTime)

	log.Println("도서관(2/3)")
	scrape(academic, "도서관)")
	time.Sleep(crawler.WaitTime)

	log.Println("도서관(3/3)")
	scrape(event, "도서관)")
	time.Sleep(crawler.WaitTime)
	log.Println("Completed scrape 도서관")
}
