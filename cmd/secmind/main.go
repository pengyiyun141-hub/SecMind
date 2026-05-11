package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	//"text/scanner"
	//"strings"
	//"github.com/PuerkitoBio/goquery"
	//"golang.org/x/text/message"
	"bytes"
	"encoding/json"
	//"io"
	"os"
)

type Article struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

func main() {
	urls_file, err:= os.Open("../../configs/urls.txt")

	if err != nil {
		log.Fatal("文件打开失败:", err)
	}

	scanner := bufio.NewScanner(urls_file)

	var urls_str []string
	for scanner.Scan() {
		line := scanner.Text()

		urls_str = append(urls_str, line)
	}
	
		var xmlData_slice []Article1
		for article := range Fetch(urls_str) {
    		xmlData_slice = append(xmlData_slice, article)
		}

	i := len(xmlData_slice)

	if  i > 0 {
		

		for _, article := range xmlData_slice {
    		fmt.Printf("标题 %d: %s\n[%s] 源:[%s]\n\n", article.Id+1, article.Title, article.Link, article.Source)
		}
	}
	fmt.Println("")
	
	SaveToMD(xmlData_slice)

}	