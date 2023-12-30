package crawl

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type CrawlVideoParam struct {
	webUrl   string
	diskPath string
}

func CrawlVideo(param CrawlVideoParam) {
	var m3u8Url string
	GetM3u8Url(param.webUrl, &m3u8Url)
	SaveHls(HlsSaveParam{
		url:      m3u8Url,
		diskPath: param.diskPath,
	})
}

type VideoElement struct {
	Url        string `json:"url"`
	Pic        string `json:"pic"`
	Type       string `json:"type"`
	Thumbnails string `json:"thumbnails"`
}

type DataConfig struct {
	Video VideoElement `json:"video"`
}

func GetM3u8Url(webUrl string, m3u8Url *string) {
	doc, err := goquery.NewDocument(webUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".dplayer").Each(func(i int, s *goquery.Selection) {
		// 获取 "myattribute" 属性的值
		value, exists := s.Attr("data-config")
		if exists {
			log.Println(value)
			var dataConfig DataConfig
			err := json.Unmarshal([]byte(value), &dataConfig)
			if err != nil {
				log.Fatal(err)
			}
			*m3u8Url = dataConfig.Video.Url
		}
	})
}

type HlsSaveParam struct {
	url      string
	diskPath string
}

func SaveHls(param HlsSaveParam) {
	resp, err := http.Get(param.url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	p, _, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		panic(err)
	}

	playlist, ok := p.(*m3u8.MediaPlaylist)
	if !ok {
		panic("Invalid playlist")
	}

	os.MkdirAll(param.diskPath, os.ModePerm)

	localPlaylist, err := m3u8.NewMediaPlaylist(uint(playlist.Count()), uint(playlist.Count()))
	if err != nil {
		panic(err)
	}

	key := playlist.Key
	if key != nil {
		resp, err := http.Get(key.URI)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		u, err := url.Parse(key.URI)
		if err != nil {
			panic(err)
		}

		keyName := filepath.Base(u.Path)
		keyFilePath := filepath.Join(param.diskPath, keyName)
		out, err := os.Create(keyFilePath)
		if err != nil {
			panic(err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}

		key.URI = keyName
		localPlaylist.Key = key
	}

	for _, v := range playlist.Segments {
		if v != nil {
			resp, err := http.Get(v.URI)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			u, err := url.Parse(v.URI)
			if err != nil {
				panic(err)
			}

			tsName := filepath.Base(u.Path)
			tsFilePath := filepath.Join(param.diskPath, tsName)
			out, err := os.Create(tsFilePath)
			if err != nil {
				panic(err)
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				panic(err)
			}

			localPlaylist.Append(tsName, v.Duration, "")
		}
	}

	buf := localPlaylist.Encode()
	if err := os.WriteFile(filepath.Join(param.diskPath, "index.m3u8"), buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
