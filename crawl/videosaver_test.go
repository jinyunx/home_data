package crawl

import (
	"testing"
)

var hlsURL string = ""
var diskPath string = "../data/103942/video"
var webUrl string = "https://hy7uz1.jrkrta.com/archives/103942/"

// 测试函数
func TestHlsSave(t *testing.T) {
	vs := VideoSaver{
		diskPath: diskPath,
		m3u8Url:  hlsURL,
		selector: ".dplayer",
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
	}
	if err := vs.Run(); err != nil {
		t.Fatal(err)
	}

	t.Log("m3u8Url:", vs.M3u8Url())
}
