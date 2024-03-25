package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var svrUrl string

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage exe serverIp")
		return
	}
	svrUrl = os.Args[1]

	httpServer()
}

func httpGet(url string) (error, []byte) {
	// 发送GET请求
	response, err := http.Get(url)
	if err != nil {
		log.Println("Error making GET request:", err)
		return err, []byte{}
	}
	defer response.Body.Close() // 确保关闭响应的Body

	// 读取响应正文
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return err, []byte{}
	}
	return nil, body
}

// localHandler 是一个HTTP处理函数，它写入一个简单的响应
func localHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("localHandler running")
	err, content := httpGet(svrUrl)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}
	origData, err := AesDecrypt(string(content), []byte(key))
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprint(w, string(origData))
	//fmt.Fprint(w, string(content))
	log.Println("localHandler running end")
}

func httpServer() {
	// 设置路由和处理函数
	http.HandleFunc("/", localHandler)

	// 定义服务器监听的端口
	port := ":7070"
	fmt.Printf("Server is running at %s\n", port)

	// 启动HTTP服务器
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
