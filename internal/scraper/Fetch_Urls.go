package scarper

import (
	"fmt"
	"log"
	"net/http"
	"secmind/internal/model"
	"secmind/internal/parser"
	"sync"
)

func Fetch(urls []string) <-chan model.Article {
	var wg sync.WaitGroup

	ch := make(chan model.Article, len(urls))

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {

			defer wg.Done()

			resp, err := http.Get(url)

			if err != nil {
				log.Printf("请求失败:[URL]: %s, %s", url, err)
				return
			}

			defer resp.Body.Close()

			fmt.Println("开始抓取：", url)

			xmlData, err := parser.Parse(resp.Body, url)

			if err != nil {
				log.Printf("解析失败:%s，%s", url, err)
			}

			for _, article := range xmlData {
				ch <- article
			}
		}(url)
	}
	go func() {
		wg.Wait() // ① 等待所有 goroutine 完成
		close(ch) // ② 所有任务完成后关闭通道
	}()
	return ch
}
