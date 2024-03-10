package crawl

import (
	"testing"
)

var hlsURL string = ""
var diskPath string = "../data/"
var webUrl string = "https://abandon.ucqhumktg.com/archives/123659/"

// 测试函数
func TestHlsSave(t *testing.T) {
	vs := VideoSaver{
		diskPath: diskPath,
		m3u8Url:  hlsURL,
		selector: ".dplayer",
		name:     "103942",
	}
	if err := vs.SaveHls(); err != nil {
		t.Fatal(err)
	}
}

func TestGetM3u8Url(t *testing.T) {
	vs := VideoSaver{
		webUrl:   webUrl,
		diskPath: diskPath,
		selector: ".dplayer",
		name:     "103942",
	}
	err, m3u8Url := vs.GetM3u8Url()
	if err != nil {
		t.Fatal(err)
	}

	if err := vs.SaveHls(); err != nil {
		t.Fatal(err)
	}

	t.Log("m3u8Url:", m3u8Url)
}

func TestCrawlVideo(t *testing.T) {
	vs := VideoSaver{
		webUrl:   webUrl,
		diskPath: diskPath,
		selector: ".dplayer",
		name:     "123659",
	}
	if err := vs.Run(); err != nil {
		t.Fatal(err)
	}

	t.Log("m3u8Url:", vs.M3u8Url())
}
