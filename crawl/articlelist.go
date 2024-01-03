package crawl

import (
	"github.com/PuerkitoBio/goquery"
	"log"
)

type ArticleList struct {
	PageUrl string
}

func (a *ArticleList) GetWebUrlList() (error, []string) {
	log.Println("GetWebUrlList running")

	doc, err := goquery.NewDocument(a.PageUrl)
	if err != nil {
		log.Println("goquery.NewDocument fail", err)
		return err, nil
	}

	var result []string
	doc.Find("#index article a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			result = append(result, href)
			log.Println(href)
		}
	})
	return nil, result
}
