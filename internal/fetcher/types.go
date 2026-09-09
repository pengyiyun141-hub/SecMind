package fetcher

import (
	"context"
	"net/http"
	"secmind/configs"
	"secmind/internal/article"
)


type Fetcher interface {
    Fetch(ctx context.Context, fetchRequest FetchRequest)(FetchResult, error)
}

type FeedFetcher struct {
	httpClient	   http.Client
	FeedSourceInfo configs.FeedSourceInfo
}

type FetchRequest struct {
	SourceInfo	*configs.SourceInfo
}

type FetchResult struct {
	FeedArticles	[]article.FeedArticle
	NewCount		int
}