package fetcher

import (
	"context"
	"secmind/configs"
	"secmind/internal/article"
)


type Fetcher interface {
    Fetch(ctx context.Context, fetchRequest FetchRequest)(FetchResult, error)
}

type FetchRequest struct {
	SourceInfo	*configs.SourceInfo
}

type FetchResult struct {
	FeedArticles	[]article.FeedArticle
	NewCount		int
}