package image

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chromedp/chromedp"
)

var ctx context.Context
var allocCancelFunc context.CancelFunc
var cancel context.CancelFunc

func init() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),              // 启用无头模式
		chromedp.Flag("no-sandbox", true),            // 禁用沙箱
		chromedp.Flag("disable-gpu", true),           // 禁用 GPU 加速
		chromedp.Flag("disable-dev-shm-usage", true), // 解决共享内存问题
	)

	ctx, allocCancelFunc = chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(ctx)

	// 创建一个通道，用于接收信号
	sigChan := make(chan os.Signal, 1)
	// 监听 SIGINT（Ctrl+C）和 SIGTERM（kill）信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 来处理信号
	go func() {
		<-sigChan // 等待信号
		allocCancelFunc()
		cancel()
		os.Exit(0) // 退出程序
	}()
}

func DarwImg(url string, jsonData string, writePath string) {
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, imageShot(url, jsonData, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(writePath, buf, 0o644); err != nil {
		log.Fatal(err)
	}
}

func imageShot(urlStr string, jsonData string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(`#dataJson`, chromedp.ByID),        // 替换为输入框的选择器
		chromedp.SendKeys(`#dataJson`, jsonData, chromedp.ByID), // 输入内容
		chromedp.Click(`#genImg`, chromedp.ByID),                // 替换为按钮的选择器
		chromedp.WaitVisible(`#res_image`, chromedp.ByID),
		chromedp.Screenshot(`#res_image`, res),
	}
}
