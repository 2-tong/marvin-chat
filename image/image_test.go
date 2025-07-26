package image

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"testing"
)

func TestImageTa(t *testing.T) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),              // 启用无头模式
		chromedp.Flag("no-sandbox", true),            // 禁用沙箱
		chromedp.Flag("disable-gpu", true),           // 禁用 GPU 加速
		chromedp.Flag("disable-dev-shm-usage", true), // 解决共享内存问题
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(`file:///D:/temp/aaa.html`, 90, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("D://temp//fullScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

func fullScreenshot(urlStr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(`#dataJson`, chromedp.ByID),                           // 替换为输入框的选择器
		chromedp.SendKeys(`#dataJson`, `{"curr":{"name":"诸王峡谷" }}`, chromedp.ByID), // 输入内容
		chromedp.Click(`#genImg`, chromedp.ByID),                                   // 替换为按钮的选择器
		chromedp.WaitVisible(`#res_image`, chromedp.ByID),
		chromedp.Screenshot(`#res_image`, res),
	}
}
