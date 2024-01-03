package crawl

import (
	"github.com/jinyunx/home_data/database"
	"github.com/jinyunx/home_data/taskqueue"
	"github.com/jinzhu/gorm"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type FetchParam struct {
	WebUrl   string
	DiskPath string
}

type FetchTask struct {
	Task *taskqueue.TaskQueue
	db   *gorm.DB
}

func NewCrawlTask(diskPath string, dbName string) *FetchTask {
	os.MkdirAll(diskPath, os.ModePerm)

	db := database.GetDB(diskPath, dbName)
	return &FetchTask{
		Task: taskqueue.NewTaskQueue(),
		db:   db,
	}
}
func (c *FetchTask) WaitToStop() {
	c.Task.WaitToStop()
}

func (c *FetchTask) ProcessOne(param FetchParam, name string) {
	log.Println("ProcessOne running")

	var dbInfo database.CrawlStatus
	dbInfo.Status = taskqueue.TaskStatusRunning
	dbInfo.WebUrl = param.WebUrl
	dbInfo.Name = name

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		s := Screenshot{
			name:     name,
			webUrl:   param.WebUrl,
			diskPath: param.DiskPath,
			timeout:  5 * time.Minute,
		}
		err := s.DoScreenshot()
		if err != nil {
			dbInfo.ScreenshotError = err.Error()
		}
		wg.Done()
	}()

	go func() {
		vs := VideoSaver{
			name:     name,
			webUrl:   param.WebUrl,
			diskPath: param.DiskPath,
			selector: ".dplayer",
		}

		err := vs.Run()
		if err != nil {
			dbInfo.VideoSaverError = err.Error()
		}
		dbInfo.M3u8Url = vs.m3u8Url
		wg.Done()
	}()

	wg.Wait()

	if len(dbInfo.ScreenshotError) > 0 || len(dbInfo.VideoSaverError) > 0 {
		dbInfo.Status = taskqueue.TaskStatusFail
	} else {
		dbInfo.Status = taskqueue.TaskStatusDone
	}
	c.db.Create(&dbInfo)
}

func (c *FetchTask) AddCrawlTask(param FetchParam) error {
	os.MkdirAll(param.DiskPath, os.ModePerm)
	u, err := url.Parse(param.WebUrl)
	if err != nil {
		log.Println("url.Parse fail", err, param.WebUrl)
		return err
	}
	name := filepath.Base(u.Path)

	c.Task.Add(name, func() int32 {
		c.ProcessOne(param, name)
		return 0
	})
	return nil
}
