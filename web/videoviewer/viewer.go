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
)

type WebData struct {
	Img  string
	M3u8 string
}

func main() {
	View("/Volumes/sata11-136XXXX0904/51cg/data")
}

func View(diskPath string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Path", r.URL.Path)

		name := filepath.Base(r.URL.Path)
		_, err := strconv.Atoi(filepath.Base(r.URL.Path))
		if err != nil {
			if name == "index.m3u8" {
				content, err := ioutil.ReadFile(filepath.Join(diskPath, r.URL.Path))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// 追加#EXT-X-ENDLIST，不然会识别为直播
				fmt.Fprintln(w, string(content))
				fmt.Fprintln(w, "#EXT-X-ENDLIST")
				log.Println(name, "#EXT-X-ENDLIST")
			} else {
				http.FileServer(http.Dir(diskPath)).ServeHTTP(w, r)
			}
			return
		} else {
			log.Println("name", name)
		}

		dir, _ := os.Getwd()
		log.Println("dir", dir)

		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Println("template.ParseFiles fail", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := WebData{
			Img:  name + ".png",
			M3u8: "video/index.m3u8",
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("tmpl.Execute fail", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
