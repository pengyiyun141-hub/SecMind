package main

import (
	"fmt"
	"log"
	"secmind/internal/analyzer"
	"secmind/internal/model"
	"secmind/internal/scraper"
	"secmind/internal/storage"
)

func main() {
	//urls_file, err := os.Open("configs/urls.txt")

	var sourceMap map[string]string
	sourceMap, err := scraper.LoadSourceMap("configs/sourceMap.json")
	if err != nil {
		log.Fatal("文件打开失败:", err)
		return
	}
	fmt.Println("sourceMap加载成功")

	var shortsource []string
	var realsource []string
	for ss, rs := range sourceMap {
		shortsource = append(shortsource, ss)
		fmt.Printf("%s:", ss)
		realsource = append(realsource, rs)
		fmt.Printf("%s\n\n", rs)
	}

	/*
			scanner := bufio.NewScanner(urls_file)
			if err := scanner.Err(); err != nil {
		    	log.Fatal("读取文件时发生错误:", err)
			}

			var urls_str []string
			for scanner.Scan() {
				line := scanner.Text()
				urls_str = append(urls_str, line)
			}
	*/

	var xmlData_slice []model.Article
	for article := range scraper.Fetch(sourceMap) {
		xmlData_slice = append(xmlData_slice, article)
	}

	l := len(xmlData_slice)

	if l > 0 {
		for _, article := range xmlData_slice {
			fmt.Printf("标题 %d: %s\n[%s] 源:[%s]\n\n", article.Id, article.Title, article.Link, article.Source)
		}
	}
	fmt.Println("")

	analyzer.AnalyzeByAI(xmlData_slice, sourceMap)
	storage.SaveToMD(xmlData_slice)

}
