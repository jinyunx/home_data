package main

import (
	"github.com/jinyunx/home_data/crawl"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

var diskPath = "data/"
var dbName = "test.db"
var host = "https://h2avz1.ilfrjq.com/"

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(pwd)

	diskPath = filepath.Join(pwd, diskPath)

	//dbviewer.View(diskPath, dbName)

	task := crawl.NewCrawlTask(diskPath, dbName)

	for i := 47; i < 1000; i += 1 {
		pageUrl, _ := url.JoinPath(host, "page/", strconv.Itoa(i))
		log.Println("pageUrl", pageUrl)
		a := crawl.ArticleList{PageUrl: pageUrl}
		err, pathList := a.GetWebUrlList()
		if err != nil {
			log.Fatal("a.GetWebUrlList", err)
		}
		for _, v := range pathList {
			webUrl, _ := url.JoinPath(host, v)
			log.Println("webUrl", webUrl)
			task.AddCrawlTask(crawl.FetchParam{
				WebUrl:   webUrl,
				DiskPath: diskPath,
			})
		}
		task.Wait()
	}
}
