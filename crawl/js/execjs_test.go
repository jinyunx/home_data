package js

import (
	"os"
	"testing"
)

func TestDecryptImage(t *testing.T) {
	content, err := os.ReadFile("2024011011580930711.jpeg")
	if err != nil {
		t.Fatal("ioutil.ReadFile fail", err)
		return
	}
	img, err := DecryptImage(content)
	if err != nil {
		t.Fatal("DecryptImage fail", err)
		return
	}
	os.WriteFile("2024011011580930711_decode.jpg", img, os.ModePerm)
}
