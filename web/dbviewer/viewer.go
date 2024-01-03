package dbviewer

import (
	"github.com/jinyunx/home_data/database"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type PageData struct {
	Entries  []database.CrawlStatus
	PrevPage int
	NextPage int
}

func View(diskPath string, dbName string) {
	db := database.GetDB(diskPath, dbName)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page <= 0 {
			page = 1
		}

		var entries []database.CrawlStatus
		result := db.Limit(20).Offset((page - 1) * 20).Find(&entries)
		if result.Error != nil {
			log.Println("db.Find fail", result.Error)
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(filepath.Join(diskPath, "web/dbviewer.html"))
		if err != nil {
			log.Println("template.ParseFiles fail", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := PageData{
			Entries:  entries,
			PrevPage: page - 1,
			NextPage: page + 1,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println("tmpl.Execute fail", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
