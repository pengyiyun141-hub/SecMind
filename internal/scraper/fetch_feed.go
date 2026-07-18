package scraper

import (
	//"encoding/json"
	"fmt"
	//"io"
	"log"
	"net/http"
	//"os"
	"secmind/internal/article"
	"secmind/internal/parser"
	"sync"
)

func Fetch(sourceMap map[string]string) <-chan article.Article {
	var wg sync.WaitGroup

	ch := make(chan article.Article)

	for shortsource, realurl := range sourceMap {
		wg.Add(1)

		go func(url string) {

			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				log.Printf("请求失败:[URL]: %s, %s", realurl, err)
				return
			}

			defer resp.Body.Close()

			fmt.Println("开始抓取：", url)

			xmlData, err := parser.Parse(resp.Body, shortsource)
			if err != nil {
				log.Printf("解析失败:%s，%s", url, err)
			}

			for _, article := range xmlData {
				ch <- article
			}
		}(realurl)
	}
	go func() {
		wg.Wait() // ① 等待所有 goroutine 完成
		close(ch) // ② 所有任务完成后关闭通道
	}()
	return ch
}

/*
func LoadSourceMap(source_file_path string) (map[string]string, error) {

	map_file, err := os.Open(source_file_path)
	if err != nil {
		log.Printf("源映射文件加载失败:%s", err)
		return nil, err
	}
	defer map_file.Close()

	map_file_data, err := io.ReadAll(map_file)
	if err != nil {
		log.Printf("从map_file中加载内容失败:%s", err)
		return nil, err
	}

	var sourceMap map[string]string
	err = json.Unmarshal(map_file_data, &sourceMap)

	return sourceMap, err
}
*/