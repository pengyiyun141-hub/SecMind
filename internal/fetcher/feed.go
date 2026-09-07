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
