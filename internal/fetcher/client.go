package fetcher

import (
	"net/http"
	"time"
	"context"
)

type FetcherClient struct {
	HttpClient	*http.Client
}

type Options struct {
	HttpTimeout time.Duration
}

func NewFetcherClient(opts Options) (*FetcherClient) {
	fetcherClinent := &FetcherClient{
		HttpClient: &http.Client{
			Timeout: opts.HttpTimeout,
		},
	}

	return fetcherClinent
}

func (FetcherClient *FetcherClient) Fetch(ctx context.Context, fetchRequest FetchRequest)(FetchResult, error) {
	switch fetchRequest.SourceInfo.Kind {
	case "feed":
		f := &FeedFetcher{
			httpClient: *FetcherClient.HttpClient,
			
		}
	}
}

