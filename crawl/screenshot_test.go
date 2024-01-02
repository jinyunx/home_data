package crawl

import (
	"os"
	"testing"
	"time"
)

// 测试函数
func TestScreenshot(t *testing.T) {
	diskPath := "../data/103942/"
	os.MkdirAll(diskPath, 0755)
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
