package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var ing = make(chan int, 0)

func main() {
	//	spider("http://book.zongheng.com/chapter/533225/10058291.html")
	//雪鹰领主
	//	cs := getmenus("http://www.boluoxs.com/biquge/0/420/")
	//	fmt.Println(cs[:10])
	//	getText("http://www.boluoxs.com/biquge/0/420/xs285541.html")
	//	Remix(cs, 10)
	//	for {
	//	}
	a := make([]int, 100)
	for k := range a {
		a[k] = k + 1
		fmt.Printf("\r%v", k+1)
	}
}

// Chapter 书籍章节
type Chapter struct {
	CID  int    //章节编号
	Name string //章节名称
	Href string //链接
	Text string //内容
}

// Content 内容
type Content struct {
	CID  int
	Text string
}

// 获取目录
func getmenus(link string) (cs []Chapter) {
	cs = make([]Chapter, 0, 1024)
	defer func() {
		fmt.Print("共有", len(cs), "章")
	}()
	doc, err := goquery.NewDocument(link)
	if err != nil {
		fmt.Println("请求出错", err)
		return
	}
	//字符编码转换
	ht, _ := doc.Html()
	reader := transform.NewReader(strings.NewReader(ht),
		simplifiedchinese.GBK.NewDecoder())
	data, _ := ioutil.ReadAll(reader)
	doc.SetHtml(string(data))

	//章节
	chapters := doc.Find("#chapterlist")
	lis := chapters.Children().Filter("li")
	for i := 0; i < lis.Length(); i++ {
		li := lis.Eq(i)
		c := Chapter{
			CID:  i,
			Name: li.Text(),
		}
		href, ok := li.Children().Attr("href")
		if ok {
			if strings.Contains(href, "www.") {
				c.Href = href
			} else {
				c.Href = link + href
			}
		}
		cs = append(cs, c)
		if i == 100 {
			return
		}
	}
	return
}

// 获取文本内容
func getText(link string) (text string) {
	doc, err := goquery.NewDocument(link)
	if err != nil {
		fmt.Println("请求出错", err)
		return
	}
	//字符编码转换
	ht, _ := doc.Html()
	reader := transform.NewReader(strings.NewReader(ht),
		simplifiedchinese.GBK.NewDecoder())
	data, _ := ioutil.ReadAll(reader)
	doc.SetHtml(string(data))

	//内容
	text = doc.Find("#book_text").Text()
	return
}

// 内容整合 n n个章节一起处理
func Remix(cs []Chapter, n int) {
	if n == 0 || len(cs) == 0 {
		return
	}
	if n > len(cs) {
		n = len(cs)
	}
	piece := 0
	if len(cs)%n == 0 {
		piece = len(cs) / n
	} else {
		piece = len(cs)/n + 1
	}
	for i := 0; i < piece; i++ {
		if i == len(cs)/n {
			go remix(cs[i*n:], strconv.Itoa(i*n+1)+"-"+strconv.Itoa(len(cs)))
			break
		}
		go remix(cs[i*n:i*n+n], strconv.Itoa(i*n+1)+"-"+strconv.Itoa(i*n+n))
	}

}

func remix(cs []Chapter, index string) {
	reader := bytes.NewBufferString("")
	f, err := os.OpenFile("C:/Users/xs/Documents/test_sp"+index+".txt",
		os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()
	for _, v := range cs {
		fmt.Print(v.CID)
		reader.WriteString(strconv.Itoa(v.CID))
		//		v.Text = getText(v.Href)
		//		reader.WriteString(v.Text)
		reader.WriteString("\n")
	}
	_, err = f.Write(reader.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\n" + index + "写入完成")
}

//进度显示
func progress(total, step int) {
	t := []rune("[" + strings.Repeat(" ", 50) + "]")
	i := 0
	for {
		a := <-ing
		if a > 0 {
			t[i] = '='
			i++
		}
	}

}

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
