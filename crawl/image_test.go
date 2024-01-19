package crawl

import "testing"

func TestGetImgUrlList(t *testing.T) {
	s := ImageSaver{
		name:     "109915",
		webUrl:   "https://h2enz2.ewkkgy.com/archives/109915/",
		diskPath: "/Users/onexie/GoProjects/home_data/data",
	}
	err, _ := s.GetImgUrlList()
	if err != nil {
		t.Fatal(err)
		return
	}
}
func TestSaveImg(t *testing.T) {
	s := ImageSaver{
		name:     "109915",
		webUrl:   "https://h2enz2.ewkkgy.com/archives/109915/",
		diskPath: "/Users/onexie/GoProjects/home_data/data",
		jsPath:   "./js",
	}
	err := s.SaveImg()
	if err != nil {
		t.Fatal(err)
		return
	}
}
