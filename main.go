package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	// 读取网页内容
	buf := make([]byte, 1024*4)
	for {
		n, err2 := resp.Body.Read(buf)
		if err2 != nil {
			if err2 == io.EOF {
				// 读取完毕
				break
			} else {
				fmt.Println("resp.Body.Read err = ", err2)
				return
			}
		}
		result += string(buf[:n])
	}
	return
}

func SpiderPage(i int, page chan<- int) {
	url := "https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn=" +
		strconv.Itoa((i-1)*50)
	fmt.Println(url)

	// 爬取网页内容
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("HttpGet err = ", err)
		return
	}

	// 把内容写入到文件
	fileName := "/Users/yangtongchun/Desktop/" + strconv.Itoa(i) + ".html"
	f, err1 := os.Create(fileName)
	if err1 != nil {
		fmt.Println("os.Create err = ", err1)
		return
	}
	f.WriteString(result)
	f.Close()
	page <- i
}

func DoWork(start, end int) {
	fmt.Println("正在爬取", start, "到", end, "的页面")
	page := make(chan int)
	// 明确爬取的范围或者网址
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d页爬取完毕\n", <-page)
	}
}

func main() {
	var start, end int
	fmt.Println("请输入起始页 (>= 1)：")
	fmt.Scan(&start)
	fmt.Println("请输入终止页 (>= 起始页)：")
	fmt.Scan(&end)

	DoWork(start, end)
}
