package crawl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type VideoSaver struct {
	webUrl   string
	diskPath string
	name     string
	selector string

	m3u8Url string
}

func (vs *VideoSaver) M3u8Url() string {
	return vs.m3u8Url
}

func (vs *VideoSaver) Run() error {
	log.Println("VideoSaver running")

	err, m3u8Url := vs.GetM3u8Url()
	if err != nil {
		return err
	}
	vs.m3u8Url = m3u8Url

	return vs.SaveHls()
}

type DataConfig struct {
	Url string `json:"url"`
}

func HttpGet(url string) (resp *http.Response, err error) {
	// 创建一个HTTP客户端
	client := &http.Client{}

	// 创建一个请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 设置User-Agent头部，模拟Chrome浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	// 发起请求
	return client.Do(req)
}

func (vs *VideoSaver) GetM3u8Url() (error, string) {
	log.Println("GetM3u8Url running")

	doc, err := goquery.NewDocument(vs.webUrl)
	if err != nil {
		log.Println("goquery.NewDocument fail", err)
		return err, ""
	}

	var m3u8Url string
	var findErr error
	doc.Find(vs.selector).Each(func(i int, s *goquery.Selection) {
		// 获取 "myattribute" 属性的值
		value, exists := s.Attr("data-config")
		if exists {
			log.Println(value)
			var dataConfig DataConfig
			err := json.Unmarshal([]byte(value), &dataConfig)
			if err != nil {
				log.Println("json.Unmarshal fail", err)
				findErr = err
				return
			}
			log.Println("dataConfig", dataConfig)
			m3u8Url = dataConfig.Url
		}
	})
	if findErr != nil {
		return findErr, ""
	}
	if len(m3u8Url) == 0 {
		errMsg := fmt.Sprintf("m3u8 not found, selector:%v", vs.selector)
		return errors.New(errMsg), ""
	}
	return nil, m3u8Url
}

func (vs *VideoSaver) SaveHls() error {
	log.Println("SaveHls running")

	filePath := filepath.Join(vs.diskPath, vs.name, "video")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Println("os.MkdirAll fail", err, filePath)
		return err
	}

	m3u8Path := filepath.Join(filePath, "index.m3u8")
	if _, err := os.Stat(m3u8Path); err == nil {
		log.Println(m3u8Path, "exists")
		return nil
	}

	resp, err := HttpGet(vs.m3u8Url)
	if err != nil {
		log.Println("http.Get fail", err, vs.m3u8Url)
		return err
	}
	defer resp.Body.Close()

	p, _, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		log.Println("m3u8.DecodeFrom fail", err, resp.Body)
		return err
	}

	playlist, ok := p.(*m3u8.MediaPlaylist)
	if !ok {
		return errors.New("Invalid playlist")
	}

	localPlaylist, err := m3u8.NewMediaPlaylist(uint(playlist.Count()), uint(playlist.Count()))
	if err != nil {
		return err
	}

	key := playlist.Key
	if key != nil {
		resp, err := HttpGet(key.URI)
		if err != nil {
			log.Println("http.Get fail", err, key.URI)
			return err
		}
		defer resp.Body.Close()

		u, err := url.Parse(key.URI)
		if err != nil {
			log.Println("url.Parse fail", err, key.URI)
			return err
		}

		keyName := filepath.Base(u.Path)
		keyFilePath := filepath.Join(filePath, keyName)
		out, err := os.Create(keyFilePath)
		if err != nil {
			log.Println("os.Create fail", err, keyFilePath)
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Println("io.Copy fail", err, keyFilePath)
			return err
		}

		key.URI = keyName
		localPlaylist.Key = key
	}

	for _, v := range playlist.Segments {
		if v != nil {
			resp, err := HttpGet(v.URI)
			if err != nil {
				log.Println("http.Get fail", err, v.URI)
				return err
			}
			defer resp.Body.Close()

			u, err := url.Parse(v.URI)
			if err != nil {
				log.Println("url.Parse fail", err, v.URI)
				return err
			}

			tsName := filepath.Base(u.Path)
			tsFilePath := filepath.Join(filePath, tsName)
			out, err := os.Create(tsFilePath)
			if err != nil {
				log.Println("os.Create fail", err, tsFilePath)
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Println("io.Copy fail", err, tsFilePath)
				return err
			}

			localPlaylist.Append(tsName, v.Duration, "")
		}
	}

	buf := localPlaylist.Encode()
	if err := os.WriteFile(m3u8Path, buf.Bytes(), 0o644); err != nil {
		log.Println("os.WriteFile fail", err, filePath)
		return err
	}
	return nil
}
