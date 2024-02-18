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
	"math/big"
	"time"
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
		return
	}

	var db = database.ConnectionDB()
	var mLatest *model.Latest
	db.Where("notice_type=?", fmt.Sprintf("%s%s", model.MainNotice, code)).Find(&mLatest)

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
			log.Println(latest.String(), articleNo.String())
			log.Println(cmp)
			continue
		} else if newLatest.Cmp(&articleNo) == -1 {
			log.Println(latest.String(), articleNo.String())
			log.Println(newLatest.String())
			newLatest.SetString(articleNoStr, 10)
		}
		embeds = append(embeds, webhook.Embed{
			Title: title,
			Url:   fmt.Sprintf("%s%s", fmt.Sprintf(target, code), href),
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
			NoticeType: fmt.Sprintf("%s%s", model.MainNotice, code),
			URL:        newLatest.String(),
		})
	} else {
		return
	}

	var subscribers []*model.Subscriber
	db.Where("Main = ?", true).Find(&subscribers)

	for _, subscriber := range subscribers {
		go crawler.Send(name, subscriber, &embeds)
	}
}

func Scrape() {
	log.Println("Start scrape 메인포털")
	log.Println("메인 포털(1/3)")
	scrape(university, "메인포털(대학공지)")
	time.Sleep(crawler.WaitTime)
	log.Println("메인 포털(2/3)")
	scrape(academic, "메인포털(학사공지)")
	time.Sleep(crawler.WaitTime)
	log.Println("메인 포털(3/3)")
	scrape(scholarship, "메인포털(장학공지)")
	time.Sleep(crawler.WaitTime)
	log.Println("Completed scrape 메인포털")
}
