package crawl

import (
	"github.com/jinyunx/home_data/taskqueue"
	"github.com/jinzhu/gorm"
	"net/url"
	"path/filepath"
	"sync"
	"time"
)

type CrawlParam struct {
	webUrl   string
	dataDisk string
}

type CrawlInfo struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex"`
	Status  int32
	Html    string
	M3u8Url string
}

type CrawlTask struct {
	Task *taskqueue.TaskQueue
}

func NewCrawlTask() *CrawlTask {
	return &CrawlTask{
		Task: taskqueue.NewTaskQueue(),
	}
}

func (c *CrawlTask) AddCrawlTask(param CrawlParam) {
	u, err := url.Parse(param.webUrl)
	if err != nil {
		panic(err)
	}

	name := filepath.Base(u.Path)
	c.Task.Add(name, func() int32 {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			DoScreenshot(ScreenshotParam{
				name:     name,
				webUrl:   param.webUrl,
				diskPath: param.dataDisk,
				timeout:  5 * time.Minute,
			})
			wg.Done()
		}()

		go func() {
			CrawlVideo(CrawlVideoParam{
				webUrl:   param.webUrl,
				diskPath: param.dataDisk,
			})
			wg.Done()
		}()

		wg.Wait()

		return 0
	})
}
