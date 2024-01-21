package crawl

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/disintegration/imaging"
	"github.com/jinyunx/home_data/crawl/js"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ImageSaver struct {
	name     string
	webUrl   string
	diskPath string
	jsPath   string
}

type TxtContent struct {
	Title string `json:"title"`
}

func (s *ImageSaver) GetImgUrlList() (error, []string) {
	log.Println("GetImgUrlList running")

	doc, err := goquery.NewDocument(s.webUrl)
	if err != nil {
		log.Println("goquery.NewDocument fail", err)
		return err, nil
	}

	maxImagCnt := 10
	var result []string
	// 查找所有有 data-xkrkllgl 属性的元素
	doc.Find("[data-xkrkllgl]").Each(func(i int, s *goquery.Selection) {
		// 获取 data-xkrkllgl 属性的值
		href, exists := s.Attr("data-xkrkllgl")
		if exists {
			// 图片太多内存吃不消
			if len(result) < maxImagCnt {
				result = append(result, href)
			}
			log.Println(href)
		}
	})
	return nil, result
}

func (s *ImageSaver) DecryptImages(imgUrls []string) (error, [][]byte) {
	log.Println("DecryptImages running")
	var result [][]byte
	for _, u := range imgUrls {
		jpgBuf, err := js.DecryptImageByUrl(u, s.jsPath)
		if err != nil {
			log.Println("js.DecryptImage fail", err)
			return err, nil
		}
		result = append(result, jpgBuf)
	}
	return nil, result
}

func (s *ImageSaver) SaveImg() error {
	log.Println("SaveImg running")
	if s.IsDone() {
		return nil
	}

	err, imgUrls := s.GetImgUrlList()
	if err != nil {
		log.Println("GetImgUrlList fail", err)
		return err
	}
	err, imgs := s.DecryptImages(imgUrls)
	if err != nil {
		log.Println("DecryptImages fail", err)
		return err
	}
	//for i, img := range imgs {
	//	os.WriteFile(strconv.Itoa(i)+".jpg", img, os.ModePerm)
	//}
	err = s.MergeImg(imgs)
	if err != nil {
		log.Println("MergeImg fail", err)
		return err
	}

	err = s.SaveTxt()
	if err != nil {
		log.Println("SaveTxt fail", err)
		return err
	}
	return nil
}

func (s *ImageSaver) SaveTxt() error {
	doc, err := goquery.NewDocument(s.webUrl) // 替换为你想要抓取的网页URL
	if err != nil {
		log.Println("NewDocument fail", err)
		return err
	}

	title := ""
	doc.Find(".post-title").Each(func(i int, s *goquery.Selection) {
		// 对于每一个匹配到的元素，获取它的文本内容
		title = s.Text()
		title = strings.Replace(title, " ", ",", -1)
	})
	content := TxtContent{Title: title}
	str, err := json.Marshal(content)
	if err != nil {
		log.Println("json.Marshal fail", err)
		return err
	}
	err = os.WriteFile(s.GetTxtPath(), str, os.ModePerm)
	if err != nil {
		log.Println("json.Marshal fail", err)
		return err
	}
	return nil
}

func (s *ImageSaver) IsDone() bool {
	p := s.GetJpgPath()
	if _, err := os.Stat(p); err == nil {
		log.Println(p, "exists")
		return true
	}
	filePath := filepath.Join(s.diskPath, s.name)
	os.MkdirAll(filePath, os.ModePerm)

	return false
}

func (s *ImageSaver) GetPngPath() string {
	filePath := filepath.Join(s.diskPath, s.name)

	return filepath.Join(filePath, s.name+".png")
}

func (s *ImageSaver) GetJpgPath() string {
	filePath := filepath.Join(s.diskPath, s.name)

	return filepath.Join(filePath, s.name+".jpg")
}

func (s *ImageSaver) GetTxtPath() string {
	filePath := filepath.Join(s.diskPath, s.name)

	return filepath.Join(filePath, s.name+".json")
}

func (s *ImageSaver) bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (s *ImageSaver) MergeImg(imgs [][]byte) error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("MergeImg start Alloc = %v MiB\n", s.bToMb(m.Alloc))

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
	runtime.ReadMemStats(&m)
	log.Printf("imgDecoded done Alloc = %v MiB\n", s.bToMb(m.Alloc))

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
		runtime.ReadMemStats(&m)
		log.Printf("longImg ready Alloc = %v MiB\n", s.bToMb(m.Alloc))

		// 创建新的图片文件
		longFile, err := os.Create(s.GetPngPath())
		if err != nil {
			log.Println("os.Create fail", err)
			return err
		}
		defer longFile.Close()

		// 将拼接后的图片编码为JPEG格式并保存到新的文件
		jpeg.Encode(longFile, longImg, &jpeg.Options{Quality: jpeg.DefaultQuality})

	}

	runtime.ReadMemStats(&m)
	log.Printf("longImg done Alloc = %v MiB\n", s.bToMb(m.Alloc))

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
		shortImg := imaging.New(width, splitHeight1+splitHeight2, image.Transparent)
		for i, img := range imgResize {
			x := (i % imgCntOneLine) * (width / imgCntOneLine)
			y := (i / imgCntOneLine) * splitHeight1
			shortImg = imaging.Paste(shortImg, img, image.Pt(x, y))
		}

		runtime.ReadMemStats(&m)
		log.Printf("shortImg ready Alloc = %v MiB\n", s.bToMb(m.Alloc))

		// 创建新的图片文件
		shortFile, err := os.Create(s.GetJpgPath())
		if err != nil {
			log.Println("os.Create fail", err)
			return err
		}
		defer shortFile.Close()

		// 将拼接后的图片编码为JPEG格式并保存到新的文件
		jpeg.Encode(shortFile, shortImg, &jpeg.Options{Quality: jpeg.DefaultQuality})

	}
	runtime.ReadMemStats(&m)
	log.Printf("shortImg done Alloc = %v MiB\n", s.bToMb(m.Alloc))
	return nil
}
