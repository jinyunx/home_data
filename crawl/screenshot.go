package crawl

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"log"
	"os"
	"time"
)

type Screenshot struct {
	name     string
	webUrl   string
	diskPath string
	timeout  time.Duration
}

func (s *Screenshot) DoScreenshot() error {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	var dataConfig string

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx,
		getTasks(s.webUrl, &buf, &dataConfig)); err != nil {
		log.Println("chromedp.Run fail", err)
		return err
	}

	imgPath := s.diskPath + s.name + ".png"
	if err := os.WriteFile(imgPath, buf, 0o644); err != nil {
		log.Println("os.WriteFile fail", err, imgPath)
		return err
	}
	log.Println("img path:", imgPath)
	log.Println("data-config:", dataConfig)
	return nil
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
