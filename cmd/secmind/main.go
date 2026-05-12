package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)


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
	
		var xmlData_slice []Article
		for article := range Fetch(urls_str) {
    		xmlData_slice = append(xmlData_slice, article)
		}

	i := len(xmlData_slice)

	if  i > 0 {
		

		for _, article := range xmlData_slice {
    		fmt.Printf("标题 %d: %s\n[%s] 源:[%s]\n\n", article.Id, article.Title, article.Link, article.Source)
		}
	}
	fmt.Println("")
	
	analyzeByAI(xmlData_slice)
	SaveToMD(xmlData_slice)

}	