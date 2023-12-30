package crawl

import (
	"testing"
)

var hlsURL string = ""
var basePath string = "../data/103942/video"

// 测试函数
func TestCrawlVideo(t *testing.T) {
	CrawlVideo(VideoParam{
		url:      hlsURL,
		diskPath: basePath,
	})
}
