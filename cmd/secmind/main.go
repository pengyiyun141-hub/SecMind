package main

import (
	"fmt"
	"log"
	"secmind/configs"
	"secmind/internal/analyzer"
	"secmind/internal/model"
	"secmind/internal/scraper"
	"secmind/internal/storage"
)

func main() {
	/*
		var sourceMap map[string]string
		sourceMap, err := scraper.LoadSourceMap("configs/sourceMap.json")
		if err != nil {
			log.Fatal("sourceMap.json文件打开失败:", err)
			return
		}
		fmt.Println("sourceMap加载成功")
	*/
	//var SecmindConfigs configs.SecmindConfigs
	fmt.Println("开始加载环境")
	SecmindConfigs, err := configs.LoadAllConfigs()
	if err != nil {
		log.Fatal("初始配置加载失败:%w", err)
	}

	//测试结构体变量存储情况。\
	fmt.Println(SecmindConfigs.Aiconfigs.Apiinfo)
	for name, model := range SecmindConfigs.Aiconfigs.Apiinfo {
		fmt.Printf("模型: %s, 温度: %s, MaxTokens: %s, 模型名：%s\n",
			name, model.Baseurl, model.Apikey, model.Modelname)
	}

	fmt.Println("环境加载成功")

	var shortsource []string
	var realsource []string
	for ss, rs := range SecmindConfigs.Feedconfigs.SouceMap {
		shortsource = append(shortsource, ss)
		realsource = append(realsource, rs)
	}

	//待封装为getFeed函数，该函数的职责为发出请求获取最新的源并返回映射和存储着信息的结构体数组。
	var xmlData_slice []model.Article
	for article := range scraper.Fetch(SecmindConfigs.Feedconfigs.SouceMap) {
		xmlData_slice = append(xmlData_slice, article)
	}

	//l := len(xmlData_slice)

	/*if l > 0 {
		for _, article := range xmlData_slice {
			fmt.Printf("标题 %d: %s\n[%s] 源:[%s]\n\n", article.Id, article.Title, article.Link, article.Source)
		}
	}*/

	articleIndex := make(map[string]*model.Article)
	for i := range xmlData_slice {
		key := fmt.Sprintf("%s-%d", xmlData_slice[i].Source, xmlData_slice[i].Id)
		articleIndex[key] = &xmlData_slice[i]
	}
	fmt.Println("")

	analyzer.AnalyzeByAI(xmlData_slice, SecmindConfigs.Feedconfigs.SouceMap, articleIndex)
	storage.SaveToMD(xmlData_slice)

}
