package fetcher

import (
	"net/http"
	"time"
	"context"
)

type FetcherClinent struct {
	HttpClient	*http.Client
}

type Options struct {
	HttpTimeout time.Duration
}

func NewFetcherClient(opts Options) (*FetcherClinent) {
	fetcherClinent := &FetcherClinent{
		HttpClient: &http.Client{
			Timeout: opts.HttpTimeout,
		},
	}

	return fetcherClinent
}

func Fetch(ctx context.Context, fetchRequest FetchRequest)(FetchResult, error) {
	
}

