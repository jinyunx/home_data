package main

import (
	"encoding/base64"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// 从网络上获取加密的二进制数据
	response, err := http.Get("https://pic.bjkkmya.cn/upload/xiao/20240110/2024011011580930711.jpeg")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 读取响应体的内容
	encryptedData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// 将二进制数据转换为 Base64 编码的字符串
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	// 创建一个新的 JavaScript 解释器
	vm := otto.New()

	// 从文件中读取解密代码
	script, err := vm.Compile("test/usr/plugins/tbxw/zzz1.js", nil)
	if err != nil {
		panic(err)
	}

	// 运行解密代码
	_, err = vm.Run(script)
	if err != nil {
		panic(err)
	}

	// 调用解密函数，传递 Base64 编码的字符串
	value, err := vm.Call("decryptImage", nil, encodedData)
	if err != nil {
		panic(err)
	}

	// 获取解密后的 Base64 编码的图片数据
	decodedImageData, err := value.ToString()
	if err != nil {
		panic(err)
	}

	// 将 Base64 编码的图片数据解码为二进制数据
	imageData, err := base64.StdEncoding.DecodeString(decodedImageData)
	if err != nil {
		panic(err)
	}

	// 将解密后的图片数据保存为 JPEG 文件
	file, err := os.Create("decrypted_image.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(imageData)
	if err != nil {
		panic(err)
	}
}
