package main

import (
	"github.com/jinyunx/home_data/crawl"
	"github.com/jinyunx/home_data/web/dbviewer"
	"log"
	"os"
	"path/filepath"
)

var diskPath = "data/"
var dbName = "test.db"

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(pwd)

	diskPath = filepath.Join(pwd, diskPath)

	dbviewer.View(diskPath, dbName)

	task := crawl.NewCrawlTask(diskPath, dbName)
	task.AddCrawlTask(crawl.FetchParam{
		WebUrl:   "https://hy85z2.xxousm.com/archives/106318/",
		DiskPath: "./data/",
	})
	task.AddCrawlTask(crawl.FetchParam{
		WebUrl:   "https://hy85z2.xxousm.com/archives/103942/",
		DiskPath: "./data/",
	})
	task.WaitToStop()
}
