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
	DoScreenshot(ScreenshotParam{
		name:     "103942",
		url:      "",
		diskPath: diskPath,
		timeout:  5 * time.Minute,
	})
}
