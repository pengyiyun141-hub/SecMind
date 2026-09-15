package fetcher

import (
	"context"
	"secmind/configs"
	"secmind/internal/article"
)


type Fetcher interface {
    Fetch(ctx context.Context, fetchRequest FetchRequest)(FetchResult, error)
}

type FeedFetcher struct {
	SourceName		string
	FetcherClient	*FetcherClient
	FeedSourceInfo	*configs.FeedSourceInfo
}

type APIFetcher struct {
	SourceName		string
	FetcherClient	*FetcherClient
	APISourceInfo	*configs.APISourceInfo
}

type FetchRequest struct {
	FromDate string
    ToDate   string
    Cursor   string
}

type FetchResult struct {
	FeedArticles	[]article.FeedArticle
	NewCount		int
}