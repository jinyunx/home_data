package crawl

import (
	"testing"
)

var hlsURL string = ""
var basePath string = "../data/103942/video"
var webUrl string = "https://hy7uz1.jrkrta.com/archives/103942/"

// 测试函数
func TestHlsSave(t *testing.T) {
	SaveHls(HlsSaveParam{
		url:      hlsURL,
		diskPath: basePath,
	})
}

func TestGetM3u8Url(t *testing.T) {
	var m3u8Url string = ""
	GetM3u8Url(webUrl, &m3u8Url)
	t.Log("m3u8Url:", m3u8Url)
}

func TestCrawlVideo(t *testing.T) {
	CrawlVideo(CrawlVideoParam{
		webUrl:   webUrl,
		diskPath: basePath,
	})
}
