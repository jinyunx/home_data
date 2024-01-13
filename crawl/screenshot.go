package crawl

import (
	"bytes"
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Screenshot struct {
	name     string
	webUrl   string
	diskPath string
	timeout  time.Duration
}

func (s *Screenshot) DoScreenshot() error {
	log.Println("DoScreenshot running", s.webUrl)
	filePath := filepath.Join(s.diskPath, s.name)
	os.MkdirAll(filePath, os.ModePerm)

	imgPath := filepath.Join(filePath, s.name+".png")
	if _, err := os.Stat(imgPath); err == nil {
		log.Println(imgPath, "exists")
		return nil
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:])//chromedp.Flag("headless", false),
	//chromedp.DisableGPU,
	//chromedp.CombinedOutput(log.Writer()),
	//chromedp.Flag("enable-logging", true),
	//chromedp.Flag("v", "1"),

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		//chromedp.WithDebugf(log.Printf),
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

	if err := os.WriteFile(imgPath, buf, 0o644); err != nil {
		log.Println("os.WriteFile fail", err, imgPath)
		return err
	}
	log.Println("img path:", imgPath)
	log.Println("data-config:", dataConfig)

	spriteImgPath := filepath.Join(filePath, s.name+".jpg")
	spriteImg(bytes.NewReader(buf), spriteImgPath)
	log.Println("sprite img path:", spriteImgPath)
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
		chromedp.Sleep(time.Second * 10),
		//chromedp.Evaluate(`document.querySelector('div.dplayer').getAttribute('data-config')`, dataConfig),
		//chromedp.FullScreenshot(res, 90),
		chromedp.Screenshot(".container", res, chromedp.NodeVisible),
	}
}

func spriteImg(r io.Reader, savePath string) {
	// 解码图片
	img, _, err := image.Decode(r)
	if err != nil {
		panic(err)
	}

	// 获取图片的边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	titleWidth := width
	titleHeigth := 425
	titleStart := 175
	titleRect := image.Rect(0, titleStart, titleWidth, titleHeigth+titleStart)
	titleImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(titleRect)

	// 计算每张切片的高度
	cutHeight := titleHeigth * 4

	// 切割图片并存储切片
	var subImgs []image.Image
	for i := 0; i < 4; i++ {
		x0 := 0
		y0 := i*cutHeight + titleHeigth + titleStart
		x1 := width
		y1 := y0 + cutHeight
		if y1 >= height {
			break
		}
		rect := image.Rect(x0, y0, x1, y1)
		subImg := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(rect)
		subImgs = append(subImgs, subImg)
	}

	// 缩放的宽高
	resizeWidth := width / 4
	resizeHeight := cutHeight / 4

	// 创建新的图片以容纳拼接后的图片
	newImgWidth := width
	newImgHeight := titleHeigth + resizeHeight
	newImg := imaging.New(newImgWidth, newImgHeight, image.Transparent)

	// 绘制第一张图片
	newImg = imaging.Paste(newImg, titleImg, image.Pt(0, 0))

	// 绘制剩余的四张图片
	for i, subImg := range subImgs {
		// 缩放图片
		scaledImg := imaging.Resize(subImg, resizeWidth, resizeHeight, imaging.Lanczos)

		// 计算绘制的位置
		pos := image.Pt(i*resizeWidth, titleHeigth)

		// 绘制缩放后的图片到新图片上
		newImg = imaging.Paste(newImg, scaledImg, pos)
	}

	// 创建新的图片文件
	newFile, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// 将拼接后的图片编码为JPEG格式并保存到新的文件
	jpeg.Encode(newFile, newImg, &jpeg.Options{Quality: jpeg.DefaultQuality})
}
