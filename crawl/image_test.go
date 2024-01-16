package crawl

import "testing"

func TestGetImgUrlList(t *testing.T) {
	err, _ := GetImgUrlList("https://h2enz2.ewkkgy.com/archives/109915/")
	if err != nil {
		t.Fatal(err)
		return
	}
}
func TestSaveImg(t *testing.T) {
	err := SaveImg("https://h2enz2.ewkkgy.com/archives/109915/")
	if err != nil {
		t.Fatal(err)
		return
	}
}
