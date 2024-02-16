package crawler

import (
	"context"
	"github.com/chromedp/chromedp"
)

func CreateCrawler() (context.Context, context.CancelFunc, context.CancelFunc) {
	allocatorContext, allocatorCancel := chromedp.NewRemoteAllocator(context.Background(), "wss://chrome:9222")
	ctx, ctxCancel := chromedp.NewContext(allocatorContext)

	return ctx, allocatorCancel, ctxCancel
}
