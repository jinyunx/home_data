package js

import (
	"os"
	"testing"
)

func TestDecryptImage(t *testing.T) {
	img, err := DecryptImageByUrl("https://pic.zhliua.cn/upload/xiao/20240104/2024010421202751563.jpeg", ".")
	if err != nil {
		t.Fatal("DecryptImage fail", err)
		return
	}
	os.WriteFile("2024011011580930711_decode.jpg", img, os.ModePerm)
}
