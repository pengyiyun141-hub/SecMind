package fetcher

import (
	"fmt"
	"net/http"
	"secmind/configs"
	"secmind/internal/article"
	"secmind/internal/parser"
)

func FetchFeed(sourceInfo *configs.SourceInfo) ([]article.FeedArticle, error) {
	resp, err := http.Get(sourceInfo.Feed.URL)
	if err != nil {
		return nil, fmt.Errorf(" [FetchFeed()]:URL请求失败:%s, %w", sourceInfo.Feed.URL, err)
	}
	defer resp.Body.Close()

	xmlData, err := parser.ParseFeed(resp.Body, sourceInfo)
	if err != nil {
		return nil, fmt.Errorf(" [FetchFeed()]:源解析失败:%s, %w", sourceInfo.SourceName, err)
	}

	return xmlData, nil
}

/*var wg sync.WaitGroup

	ch := make(chan article.FeedArticle, 10)

	for shortsource, realurl := range sourceMap {
		wg.Add(1)

		go func(url string) {

			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				log.Printf("请求失败:[URL]: %s, %s", realurl.URL, err)
				return
			}

			defer resp.Body.Close()

			fmt.Println("开始抓取：", url)

			xmlData, err := parser.ParseFeed(resp.Body, shortsource)
			if err != nil {
				log.Printf("解析失败:%s，%s", url, err)
			}

			for _, article := range xmlData {
				ch <- article
			}
		}(realurl.URL)
	}
	go func() {
		wg.Wait() // 等待所有 goroutine 完成
		close(ch) // 所有任务完成后关闭通道
	}()
	return ch
}
*/
