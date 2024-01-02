package main

import "github.com/jinyunx/home_data/crawl"

func main() {
	task := crawl.NewCrawlTask()
	task.AddCrawlTask(crawl.FetchParam{
		WebUrl:   "https://hy85z2.xxousm.com/archives/106318/",
		DiskPath: "./data/",
	})
}
