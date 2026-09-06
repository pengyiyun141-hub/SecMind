package fetcher

import (
	"context"
	"secmind/configs"
	"secmind/internal/article"
)

type FetchRequest struct {
	SourceInfo	*configs.SourceInfo
}

type FetchResult struct {
	FeedArticles	[]article.FeedArticle
	NewCount		int
}

type Fetcher interface {
    Fetch(ctx context.Context, fetchRequest FetchRequest)(FetchResult, error)
}