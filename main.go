package main

import (
	"fmt"
	"github.com/jinyunx/home_data/crawl"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var diskPath = "data/"
var dbName = "test.db"
var muRunBatchTask sync.Mutex
var pageUrl string
var errMsg string
var dirCnt int
var noImgCnt int
var noVideoCnt int

type SingleTask struct {
	TaskUrl string
}

type BatchTask struct {
	CategoryUrl string
	DetailHost  string
	StartPage   int
	EndPage     int
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(pwd)

	diskPath = filepath.Join(pwd, diskPath)
	task := crawl.NewCrawlTask(diskPath, dbName)
	go HttpSvr(task)
	go UpdateDir(diskPath)
	select {}
}

func HttpSvr(task *crawl.FetchTask) {
	http.HandleFunc("/main.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/single_task", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Path", r.URL.Path)
		s := SingleTask{
			TaskUrl: r.FormValue("TaskUrl"),
		}
		log.Println("SingleTask", s)
		RunSingleTask(task, &s)
	})

	http.HandleFunc("/batch_task", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Path", r.URL.Path)

		startPage, err := strconv.Atoi(r.FormValue("StartPage"))
		if err != nil {
			log.Println("strconv.Atoi", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		endPage, err := strconv.Atoi(r.FormValue("EndPage"))
		if err != nil {
			log.Println("strconv.Atoi", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		b := BatchTask{
			CategoryUrl: r.FormValue("CategoryUrl"),
			DetailHost:  r.FormValue("DetailHost"),
			StartPage:   startPage,
			EndPage:     endPage,
		}
		log.Println("BatchTask", b)
		if muRunBatchTask.TryLock() == true {
			go RunBatchTask(task, &b)
		} else {
			http.Error(w, "muRunBatchTask.TryLock fail", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/get_console", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Path", r.URL.Path)
		content := GetConsoleContent(task)
		fmt.Fprintln(w, content)
	})

	http.ListenAndServe(":9090", nil)
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func UpdateDir(diskPath string) {
	for true {
		t := time.Now()
		var err error
		dirCache, err := os.ReadDir(diskPath)
		if err != nil {
			log.Println("os.ReadDir fail", diskPath)
		}
		log.Println("time cost", time.Since(t))
		dirCnt = len(dirCache)
		for _, d := range dirCache {
			if isNumber(d.Name()) == false {
				continue
			}
			img := filepath.Join(diskPath, d.Name(), d.Name()+".jpg")
			if _, err := os.Stat(img); err == nil {
				noImgCnt++
			}

			video := filepath.Join(diskPath, d.Name(), "video/index.m3u8")
			if _, err := os.Stat(video); err == nil {
				noVideoCnt++
			}
		}
		time.Sleep(time.Hour)
	}
}

func GetConsoleContent(task *crawl.FetchTask) string {
	content := "dirCnt:" + strconv.Itoa(dirCnt) + "\n"
	content += "pageUrl:" + pageUrl + "\n"
	content += "======================\n\n"
	content += "ID,\tStatus,\tTimeAdd\n"

	taskList := task.PeekTask()
	for _, v := range taskList {
		status := "Running"
		if v.Status == 0 {
			status = "NotStart"
		}
		content += fmt.Sprintf("%v,\t%v,\t%v\n", v.ID, status, v.TimeAdd.Format("2006-01-02 15:04:05"))
	}
	content += "======================\n\n"
	content += "errMsg\n" + errMsg + "\n"
	content += "======================"
	return content
}

func RunSingleTask(task *crawl.FetchTask, s *SingleTask) {
	log.Println("webUrl", s.TaskUrl)
	task.AddCrawlTask(crawl.FetchParam{
		WebUrl:   s.TaskUrl,
		DiskPath: diskPath,
	})
}

func RunBatchTask(task *crawl.FetchTask, b *BatchTask) {
	for i := b.StartPage; i < b.EndPage; i += 1 {
		pageUrl, _ = url.JoinPath(b.CategoryUrl, strconv.Itoa(i)+"/")
		log.Println("pageUrl", pageUrl)
		a := crawl.ArticleList{PageUrl: pageUrl}
		err, pathList := a.GetWebUrlList()
		if err != nil {
			log.Println("a.GetWebUrlList", err)
			errMsg += "GetWebUrlList:" + err.Error() + "\n"
			continue
		}
		for _, v := range pathList {
			webUrl, _ := url.JoinPath(b.DetailHost, v)
			log.Println("webUrl", webUrl)
			task.AddCrawlTask(crawl.FetchParam{
				WebUrl:   webUrl,
				DiskPath: diskPath,
			})
		}
		task.Wait()
	}
	muRunBatchTask.Unlock()
}
