package crawl

import (
	"os"
	"testing"
	"time"
)

// 测试函数
func TestScreenshot(t *testing.T) {
	diskPath := "../data/"
	os.MkdirAll(diskPath, os.ModePerm)
	s := Screenshot{
		name:     "103942",
		webUrl:   "",
		diskPath: diskPath,
		timeout:  5 * time.Minute,
	}
	if err := s.DoScreenshot(); err != nil {
		t.Fatal(err)
	}
}
