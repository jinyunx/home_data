package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type WebData struct {
	Img  string
	M3u8 string
}

type Article struct {
	DetailRef string
	Img       string
}

type MenuData struct {
	Menu     []Article
	PrevPage int
	NextPage int
}

var dirCache []os.DirEntry

func main() {
	dirPath := "/Volumes/sata11-136XXXX0904/51cg/data"
	go UpdateDir(dirPath)
	View(dirPath)
}

func UpdateDir(diskPath string) {
	for true {
		t := time.Now()
		var err error
		dirCache, err = os.ReadDir(diskPath)
		if err != nil {
			log.Println("os.ReadDir fail", diskPath)
		}
		log.Println("time cost", time.Since(t))
		time.Sleep(time.Hour)
	}
}

func View(diskPath string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Path", r.URL.Path)

		name := filepath.Base(r.URL.Path)
		_, err := strconv.Atoi(filepath.Base(r.URL.Path))
		if err != nil {
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
			fetchDetail(name, w, r)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func fetchMenu(diskPath string, w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 0 {
		page = 0
	}

	pageSize := 5
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
		menuData.Menu = append(menuData.Menu, Article{
			DetailRef: e.Name(),
			Img:       e.Name() + "/" + e.Name() + ".png",
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

func fetchDetail(name string, w http.ResponseWriter, r *http.Request) {

	dir, _ := os.Getwd()
	log.Println("dir", dir)

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println("template.ParseFiles fail", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := WebData{
		Img:  name + "/" + name + ".png",
		M3u8: name + "/" + "video/index.m3u8",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println("tmpl.Execute fail", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
