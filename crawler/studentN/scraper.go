package studentN

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
	"net/url"
	"time"
)

const name = "학생처(대학일자리플러스센터)"
const target = "https://homepage.sch.ac.kr/student/06/01.jsp"

func scrape() {
	ctx, allocatorCancel, ctxCancel, timeoutCancel := crawler.CreateCrawler()

	defer allocatorCancel()
	defer ctxCancel()
	defer timeoutCancel()

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(target),
		chromedp.WaitVisible("#contentBody > div.jwxe_root.jwxe_board > table > tbody", chromedp.ByQuery),
		chromedp.Nodes("#contentBody > div.jwxe_root.jwxe_board > table > tbody > tr", &nodes, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Println(err)
		return
	}

	var db = database.ConnectionDB()
	var mLatest *model.Latest
	db.Where("notice_type=?", model.StudentNotice).Find(&mLatest)

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
	var date string
	var embeds []webhook.Embed

	for i, node := range nodes {
		if i == 0 {
			continue
		}
		err = chromedp.Run(ctx,
			chromedp.Text("a", &title, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("a", "href", &href, &ok, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text("td:nth-child(4)", &date, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if err != nil {
			log.Println(err)
			continue
		}
		if date != nowDate {
			continue
		}
		parsedURL, err := url.Parse(href)
		if err != nil {
			log.Println(err)
			continue
		}
		articleNoStr = parsedURL.Query().Get("article_no")
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
					Value: "학생지원팀",
				},
			},
		})
	}

	if newLatest.String() != "0" {
		db.Save(&model.Latest{
			NoticeType: model.StudentNotice,
			URL:        newLatest.String(),
		})
	} else {
		return
	}

	var subscribers []*model.Subscriber
	db.Where("Student = ?", true).Find(&subscribers)

	for _, subscriber := range subscribers {
		go crawler.Send(name, subscriber, &embeds)
	}
}

func Scrape() {
	log.Println("Start scrape 학생처(대학일자리플러스센터)")
	scrape()
	time.Sleep(crawler.WaitTime)
	log.Println("Completed scrape 학생처(대학일자리플러스센터)")
}
