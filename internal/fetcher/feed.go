package fetcher

import (
	"context"
	"fmt"
	"net/http"
	//"secmind/internal/article"
	"secmind/internal/parser"
)

func (FeedFetcher *FeedFetcher) Fetch(ctx context.Context, fetchReq FetchRequest) (FetchResult, error) {

	url := FeedFetcher.FeedSourceInfo.URL
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return FetchResult{}, err
	}

	req.Header.Set("User-Agent", "SecMind/0.1")

	resp, err := FeedFetcher.FetcherClient.HttpClient.Do(req)
	if err != nil {
		return FetchResult{}, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FetchResult{}, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	/*
		if err != nil {
			return nil, fmt.Errorf(" [FetchFeed()]:URL请求失败:%s, %w", FeedFetcher.FeedSourceInfo.URL, err)
		}
		defer resp.Body.Close()
	*/

	xmlData, err := parser.ParseFeed(resp.Body, FeedFetcher.SourceName)
	if err != nil {
		return FetchResult{}, fmt.Errorf(" [FetchFeed()]:源解析失败:%s, %w", FeedFetcher.FeedSourceInfo.URL, err)
	}

	var fetchResult FetchResult
	fetchResult = FetchResult{
		FeedArticles: xmlData,
		NewCount: len(xmlData),
	}

	return fetchResult, nil
}
