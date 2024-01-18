package js

import (
	b64 "encoding/base64"
	"log"
	"os/exec"
)

func DecryptImage(input []byte) ([]byte, error) {
	sEnc := b64.StdEncoding.EncodeToString(input)
	log.Println("imageBase64 len:", len(sEnc))

	jsFile := "crawl/js/decrypt.js"
	log.Println("jsFile:", jsFile)

	// 调用 Node.js 执行 JavaScript 文件
	cmd := exec.Command("node", jsFile, sEnc)
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return b64.StdEncoding.DecodeString(string(output))
}
