package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

func main() {
	httpSvr()
}

// helloHandler 是一个HTTP处理函数，它写入一个简单的响应
func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("helloHandler running")
	err, content := saveHtml()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}

	ciphertext, err := AesEncrypt([]byte(content), []byte(key))
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprint(w, ciphertext)
	log.Println("helloHandler end")
}

func httpSvr() {
	// 设置路由和处理函数
	http.HandleFunc("/", helloHandler)

	// 定义服务器监听的端口
	port := ":80"
	fmt.Printf("Server is running at %s\n", port)

	// 启动HTTP服务器
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func saveHtml() (error, string) {
	log.Println("saveHtml running")
	// 创建浏览器上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 用于存储页面HTML的变量
	var pageHTML string

	// 执行任务：导航到页面并获取整个页面的HTML
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.51cg.fun`), // 替换为你要访问的URL
		//chromedp.Navigate(`https://www.baidu.com`), // 替换为你要访问的URL
		chromedp.OuterHTML("html", &pageHTML), // 获取整个HTML标签的内容
	)
	if err != nil {
		log.Println(err)
		return err, ""
	}
	return nil, pageHTML
}
