package crawl

import (
	"github.com/grafov/m3u8"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type VideoParam struct {
	url      string
	diskPath string
}

func CrawlVideo(param VideoParam) {
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
