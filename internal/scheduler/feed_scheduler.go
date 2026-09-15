package scheduler

import (
	"context"
	"log/slog"
	"secmind/configs"
	"secmind/internal/article"
	"secmind/internal/fetcher"
	"time"
)

type FeedScheduler struct {
	SourceInfoMap 	map[string]*configs.SourceInfo
	FetcherClient	*fetcher.FetcherClient
	FeedTitlePool 	*article.FeedTitlePool
	BaseScheduler 	BaseScheduler
}

func NewFeedScheduler(SourceInfoMap map[string]*configs.SourceInfo, SignalCtx context.Context, FeedTitlePool *article.FeedTitlePool, FetcherClient *fetcher.FetcherClient) *FeedScheduler {
	feedScheduler := &FeedScheduler{
		SourceInfoMap: SourceInfoMap,
		FetcherClient: FetcherClient,
		FeedTitlePool: FeedTitlePool,
		BaseScheduler: *NewBaseScheduler(SignalCtx),
	}

	return feedScheduler
}

func (FeedScheduler *FeedScheduler) Start() {
	for _, sourceInfo := range FeedScheduler.SourceInfoMap {
		FeedScheduler.BaseScheduler.baseSchedulWg.Add(1)

		var f fetcher.Fetcher
		f = FeedScheduler.FetcherClient.NewFetcher(sourceInfo)

		go FeedScheduler.sourceLoop(sourceInfo, f)
	}

}

func (SourceScheduler *FeedScheduler) sourceLoop(sourceInfo *configs.SourceInfo, f fetcher.Fetcher) {
	defer SourceScheduler.BaseScheduler.baseSchedulWg.Done()

	slog.Info("[sourceLoop()]:源调度器已启动：", "source", sourceInfo.SourceName)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fReq := &fetcher.FetchRequest{}

	for {
		select {
		case <-ticker.C:
			fetchReult, err := f.Fetch(SourceScheduler.BaseScheduler.SignalCtx, *fReq)
			if err != nil {
				slog.Error("[sourceLoop()]:源返回文章失败：", "source", sourceInfo.SourceName, "error", err)
				continue
			}

			SourceScheduler.FeedTitlePool.Process(sourceInfo, fetchReult.FeedArticles)

		case <-SourceScheduler.BaseScheduler.SignalCtx.Done():
			return
		}
	}
}

func (SourceScheduler *FeedScheduler) Close() {
	SourceScheduler.BaseScheduler.Close()
}
