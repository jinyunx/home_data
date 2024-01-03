package database

import (
	"github.com/jinzhu/gorm"
	"os"
	"path/filepath"
)

type CrawlStatus struct {
	gorm.Model
	Name            string `gorm:"uniqueIndex"`
	Status          int32
	WebUrl          string
	M3u8Url         string
	ScreenshotError string
	VideoSaverError string
}

func GetDB(diskPath string, dbName string) *gorm.DB {
	os.MkdirAll(diskPath, os.ModePerm)

	dbPath := filepath.Join(diskPath, dbName)
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&CrawlStatus{})
	return db
}
