package fetcher

import (
	"net/http"
	"secmind/configs"
	"time"
)

type FetcherClient struct {
	HttpClient *http.Client
}

type Options struct {
	HttpTimeout time.Duration
}

func NewFetcherClient(opts Options) *FetcherClient {
	fetcherClinent := &FetcherClient{
		HttpClient: &http.Client{
			Timeout: opts.HttpTimeout,
		},
	}

	return fetcherClinent
}

func (FetcherClient *FetcherClient) NewFetcher(sourceInfo *configs.SourceInfo) Fetcher {
	switch sourceInfo.Kind {
	case "feed":
		f := &FeedFetcher{
			SourceName: sourceInfo.SourceName,
			FetcherClient:  FetcherClient,
			FeedSourceInfo: sourceInfo.Feed,
		}
		return f
		/*
			case "openalex":
				f := &APIFetcher{
					FetcherClient: FetcherClient,
					APISourceInfo: sourceInfo.API,
				}
				return f*/
	}
	return nil
}
