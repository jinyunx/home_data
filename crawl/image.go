package crawl

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinyunx/home_data/crawl/js"
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
	return nil
}
