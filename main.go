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
)

var diskPath = "data/"
var dbName = "test.db"
var muRunBatchTask sync.Mutex
var pageUrl string
var errMsg string

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

func GetConsoleContent(task *crawl.FetchTask) string {
	content := "pageUrl\n" + pageUrl + "\n"
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
