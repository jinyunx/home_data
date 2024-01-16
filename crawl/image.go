package crawl

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/disintegration/imaging"
	"github.com/jinyunx/home_data/crawl/js"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetImgUrlList(detailUrl string) (error, []string) {
	log.Println("GetImgUrlList running")

	doc, err := goquery.NewDocument(detailUrl)
	if err != nil {
		log.Println("goquery.NewDocument fail", err)
		return err, nil
	}

	var result []string
	// 查找所有有 data-xkrkllgl 属性的元素
	doc.Find("[data-xkrkllgl]").Each(func(i int, s *goquery.Selection) {
		// 获取 data-xkrkllgl 属性的值
		href, exists := s.Attr("data-xkrkllgl")
		if exists {
			result = append(result, href)
			log.Println(href)
		}
	})
	return nil, result
}

func DecryptImages(imgUrls []string) (error, [][]byte) {
	log.Println("DecryptImages running")
	var result [][]byte
	for _, u := range imgUrls {
		resp, err := http.Get(u)
		if err != nil {
			fmt.Println("http.Get fail", err)
			return err, nil
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		jpgBuf, err := js.DecryptImage(body)
		if err != nil {
			fmt.Println("js.DecryptImage fail", err)
			return err, nil
		}
		result = append(result, jpgBuf)
	}
	return nil, result
}

func SaveImg(detailUrl string) error {
	log.Println("SaveImg running")
	err, imgUrls := GetImgUrlList(detailUrl)
	if err != nil {
		log.Println("GetImgUrlList fail", err)
		return err
	}
	err, imgs := DecryptImages(imgUrls)
	if err != nil {
		log.Println("DecryptImages fail", err)
		return err
	}
	for i, img := range imgs {
		os.WriteFile(strconv.Itoa(i)+".jpg", img, os.ModePerm)
	}
	MergeImg(imgs)
	return nil
}

func MergeImg(imgs [][]byte) error {
	var imgDecoded []image.Image
	height := 0
	width := 800
	for _, img := range imgs {
		item, _, err := image.Decode(bytes.NewReader(img))
		if err != nil {
			log.Println()
			return err
		}
		height += item.Bounds().Dy()
		imgDecoded = append(imgDecoded, item)
	}

	// 长图
	if true {
		longImg := imaging.New(width, height, image.Transparent)
		y := 0
		for _, img := range imgDecoded {
			if img.Bounds().Dx() != width {
				img = imaging.Resize(img, width, 0, imaging.Lanczos)
			}

			longImg = imaging.Paste(longImg, img, image.Pt(0, y))
			y += img.Bounds().Dy()
		}
		// 创建新的图片文件
		longFile, err := os.Create("longImg.jpg")
		if err != nil {
			log.Println("os.Create fail", err)
			return err
		}
		defer longFile.Close()

		// 将拼接后的图片编码为JPEG格式并保存到新的文件
		jpeg.Encode(longFile, longImg, &jpeg.Options{Quality: jpeg.DefaultQuality})

	}

	if true {
		imgCnt := len(imgDecoded)
		if imgCnt > 8 {
			imgCnt = 8
		}

		numLine := 2
		imgCntOneLine := imgCnt / 2
		if imgCnt < 4 {
			imgCntOneLine = imgCnt
			numLine = 1
		}

		splitHeight1 := 0
		splitHeight2 := 0
		var imgResize []image.Image
		for i, img := range imgDecoded {
			if i > numLine*imgCntOneLine {
				break
			}
			img = imaging.Resize(img, width/imgCntOneLine, 0, imaging.Lanczos)
			imgResize = append(imgResize, img)

			if img.Bounds().Dy() > splitHeight1 && len(imgResize) <= imgCntOneLine {
				splitHeight1 = img.Bounds().Dy()
			}
			if img.Bounds().Dy() > splitHeight2 && len(imgResize) > imgCntOneLine {
				splitHeight2 = img.Bounds().Dy()
			}
		}
		shortImg := imaging.New(width, splitHeight1+splitHeight1, image.Transparent)
		for i, img := range imgResize {
			x := (i % imgCntOneLine) * (width / imgCntOneLine)
			y := (i / imgCntOneLine) * splitHeight1
			shortImg = imaging.Paste(shortImg, img, image.Pt(x, y))
		}

		// 创建新的图片文件
		shortFile, err := os.Create("short.jpg")
		if err != nil {
			log.Println("os.Create fail", err)
			return err
		}
		defer shortFile.Close()

		// 将拼接后的图片编码为JPEG格式并保存到新的文件
		jpeg.Encode(shortFile, shortImg, &jpeg.Options{Quality: jpeg.DefaultQuality})

	}
	return nil
}
