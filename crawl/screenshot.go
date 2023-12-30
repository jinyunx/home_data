package crawl

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"log"
	"os"
	"time"
)

type ScreenshotParam struct {
	name     string
	url      string
	diskPath string
	timeout  time.Duration
}

func DoScreenshot(param ScreenshotParam) {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, param.timeout)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	var dataConfig string

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx,
		getTasks(param.url, &buf, &dataConfig)); err != nil {
		log.Fatal(err)
	}

	imgPath := param.diskPath + param.name + ".png"
	if err := os.WriteFile(imgPath, buf, 0o644); err != nil {
		log.Fatal(err)
	}
	log.Println("img path:", imgPath)
	log.Println("data-config:", dataConfig)
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func getTasks(urlstr string, res *[]byte, dataConfig *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Emulate(device.IPhone12Pro),
		chromedp.Navigate(urlstr),
		//chromedp.Evaluate(`document.querySelector('div.dplayer').getAttribute('data-config')`, dataConfig),
		chromedp.FullScreenshot(res, 90),
	}
}
