package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func spider(link string) {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		fmt.Println("请求出错", err)
		return
	}
	//章节名
	name := doc.Find("h1").Text()
	fmt.Println(name)
	sele := doc.Find("#nextChapterButton")
	content := doc.Find("#chapterContent")
	//正文是否含有p标签
	child := content.Children()
	reader := bytes.NewBufferString(name + "\n")
	f, _ := os.OpenFile("d:/test_sp.txt", os.O_CREATE|os.O_APPEND, os.ModePerm)
	defer f.Close()
	for i := 0; i < child.Length(); i++ {
		reader.WriteString(child.Eq(i).Text())
		reader.WriteString("\n")
	}
	f.Write(reader.Bytes())
	next, ok := sele.Attr("href")
	if !ok {
		return
	}
	spider(next)
}

func main() {
	spider("http://book.zongheng.com/chapter/533225/10058291.html")
	//xylz
	getmenus("http://www.boluoxs.com/biquge/0/420/")
}

// 获取目录
func getmenus(link string) {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		fmt.Println("请求出错", err)
		return
	}
	//章节名
	name := doc.Find("h1").Text()
	fmt.Println(name)
	sele := doc.Find("#nextChapterButton")
}
