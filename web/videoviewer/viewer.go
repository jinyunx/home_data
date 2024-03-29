package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinyunx/home_data/crawl"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type WebData struct {
	Img   string
	M3u8  string
	Title string
}

type Article struct {
	DetailRef string
	Img       string
	Title     string
}

type MenuData struct {
	Menu     []Article
	PrevPage int
	NextPage int
}

var dirCache FileInfoSlice

func main() {
	dirPath := "../../../../data"
	go UpdateDir(dirPath)
	View(dirPath)
}

// FileInfoSlice 用于实现 sort.Interface 以便按修改时间排序
type FileInfoSlice []os.DirEntry

func (fis FileInfoSlice) Len() int {
	return len(fis)
}

func (fis FileInfoSlice) Less(i, j int) bool {
	// 使用 After 方法来判断时间的先后，这里我们按照时间逆序排序
	infoi, _ := fis[i].Info()
	infoj, _ := fis[j].Info()

	return infoi.ModTime().After(infoj.ModTime())
}

func (fis FileInfoSlice) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func UpdateDir(diskPath string) {
	for true {
		t := time.Now()
		var err error
		dirCache, err = os.ReadDir(diskPath)
		if err != nil {
			log.Println("os.ReadDir fail", diskPath)
		}
		// 按修改时间逆序排序
		sort.Sort(dirCache)
		log.Println("time cost", time.Since(t))
		time.Sleep(time.Hour * 3)
	}
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func View(diskPath string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Path", r.URL.Path)

		name := filepath.Base(r.URL.Path)
		if !isNumber(filepath.Base(r.URL.Path)) {
			if name == "index.m3u8" {
				fetchM3u8(diskPath, name, w, r)
			} else if r.URL.Path == "/" {
				fetchMenu(diskPath, w, r)
			} else {
				http.FileServer(http.Dir(diskPath)).ServeHTTP(w, r)
			}
			return
		} else {
			log.Println("name", name)
			fetchDetail(diskPath, name, w, r)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func getTitle(diskPath string, name string) string {
	p := name + "/" + name + ".json"
	file, _ := ioutil.ReadFile(filepath.Join(diskPath, p))

	data := crawl.TxtContent{}
	_ = json.Unmarshal([]byte(file), &data)
	return data.Title
}

func fetchMenu(diskPath string, w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 0 {
		page = 0
	}

	pageSize := 10
	start := page * pageSize
	end := start + pageSize
	if end >= len(dirCache) {
		end = len(dirCache)
	}
	if start >= len(dirCache) {
		log.Println("EOF", start, len(dirCache))
		http.Error(w, "EOF", http.StatusInternalServerError)
		return
	}

	var menuData MenuData
	for _, e := range dirCache[start:end] {
		if !isNumber(e.Name()) {
			log.Println("!isNumber", e.Name())
			continue
		}
		menuData.Menu = append(menuData.Menu, Article{
			DetailRef: e.Name(),
			Img:       e.Name() + "/" + e.Name() + ".jpg",
			Title:     getTitle(diskPath, e.Name()),
		})
	}
	menuData.NextPage = page + 1
	menuData.PrevPage = page - 1

	tmpl, err := template.ParseFiles("menu.html")
	if err != nil {
		log.Println("template.ParseFiles fail", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, menuData); err != nil {
		log.Println("tmpl.Execute fail", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchM3u8(diskPath string, name string, w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile(filepath.Join(diskPath, r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 追加#EXT-X-ENDLIST，不然会识别为直播
	fmt.Fprintln(w, string(content))
	fmt.Fprintln(w, "#EXT-X-ENDLIST")
	log.Println(name, "#EXT-X-ENDLIST")
}

func fetchDetail(diskPath string, name string, w http.ResponseWriter, r *http.Request) {

	dir, _ := os.Getwd()
	log.Println("dir", dir)

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println("template.ParseFiles fail", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := WebData{
		Img:   name + "/" + name + ".png",
		M3u8:  name + "/" + "video/index.m3u8",
		Title: getTitle(diskPath, name),
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("tmpl.Execute fail", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
