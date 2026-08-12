package main

import (
	"fmt"
	"log"
	"secmind/configs"
	"secmind/internal/article"
	"secmind/internal/scraper"
)

func main() {
	fmt.Println("开始加载环境")
	SecmindConfigs, err := configs.LoadAllConfigs()
	if err != nil {
		log.Fatalf("初始配置加载失败：%v", err)
	}

	//测试结构体变量存储情况。
	
	//fmt.Println(SecmindConfigs.Aiconfigs.Modelinfo["filter-title"])
	//fmt.Println(SecmindConfigs.Aiconfigs.Modelinfo["summarize-article"])
	/*for name, model := range SecmindConfigs.Aiconfigs.Modelinfo {
		fmt.Printf("模型: %s, 温度: %s, MaxTokens: %s, 模型名：%s, %s\n",
			name, model.APIKey, model.BaseURL, model.PromptSysText, model.UserPrompt)
	}*/
	/*
	var shortsource []string
	var realsource []string
	for ss, rs := range SecmindConfigs.Feedconfigs.SouceMap {
		shortsource = append(shortsource, ss)
		realsource = append(realsource, rs)
	}*/

	//待封装为getFeed函数，该函数的职责为发出请求获取最新的源并返回映射和存储着信息的结构体数组。
	/*var xmlData_slice []article.Article
	for article := range scraper.Fetch(SecmindConfigs.Feedconfigs.SouceMap) {
		xmlData_slice = append(xmlData_slice, article)
	}*/
	article.TitlePool(scraper.FetchFeed(SecmindConfigs.Feedconfigs.SouceMap))
}
