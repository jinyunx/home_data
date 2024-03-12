package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	var url string
	var outputPath string
	var quality int

	// 设置命令行参数
	flag.StringVar(&url, "url", "https://www.example.com", "the URL to take a full page screenshot of")
	flag.StringVar(&outputPath, "output", "fullpage_screenshot.png", "the output file")
	flag.IntVar(&quality, "quality", 90, "the quality of the screenshot in percentage")
	flag.Parse()

	// 创建上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 截图的二进制数据
	var buf []byte

	// 运行任务
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.FullScreenshot(&buf, quality),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 将截图保存到文件
	err = ioutil.WriteFile(outputPath, buf, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Full page screenshot saved to: %s\n", outputPath)
}
