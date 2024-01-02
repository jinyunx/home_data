package crawl

import (
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

type SqliteInfo struct {
	gorm.Model
	Name            string `gorm:"uniqueIndex"`
	Status          int32
	WebUrl          string
	M3u8Url         string
	ScreenshotError string
	VideoSaverError string
}

type FetchTask struct {
	Task *taskqueue.TaskQueue
	db   *gorm.DB
}

func NewCrawlTask() *FetchTask {
	db, err := gorm.Open("sqlite3", "./data/test.db")
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&SqliteInfo{})

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

	var dbInfo SqliteInfo
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
		dbInfo.ScreenshotError = err.Error()
		wg.Done()
	}()

	go func() {
		vs := VideoSaver{
			webUrl:   param.WebUrl,
			diskPath: param.DiskPath,
			selector: ".dplayer",
		}

		err := vs.Run()
		dbInfo.VideoSaverError = err.Error()
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
	os.MkdirAll(param.DiskPath, 0755)
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
