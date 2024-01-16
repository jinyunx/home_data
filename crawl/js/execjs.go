package js

import (
	b64 "encoding/base64"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
)

func DecryptImage(input []byte) ([]byte, error) {
	sEnc := b64.StdEncoding.EncodeToString(input)

	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)

	// 获取当前文件的目录
	dir := filepath.Dir(filename)
	jsFile := filepath.Join(dir, "decrypt.js")

	// 调用 Node.js 执行 JavaScript 文件
	cmd := exec.Command("node", jsFile, sEnc)
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return b64.StdEncoding.DecodeString(string(output))
}
