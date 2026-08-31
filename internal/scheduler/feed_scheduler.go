package scheduler

import (
	"context"
	"log"
	"log/slog"
	"secmind/configs"
	"secmind/internal/article"
	"secmind/internal/scraper"
	"time"
)

type FeedScheduler struct {
	SourceInfoMap map[string]*configs.SourceInfo
	FeedTitlePool *article.FeedTitlePool
	BaseScheduler BaseScheduler
}

func NewFeedScheduler(SourceInfoMap map[string]*configs.SourceInfo, SignalCtx context.Context, FeedTitlePool *article.FeedTitlePool) *FeedScheduler {
	feedScheduler := &FeedScheduler{
		SourceInfoMap: SourceInfoMap,
		FeedTitlePool: FeedTitlePool,
		BaseScheduler: *NewBaseScheduler(SignalCtx),
	}

	return feedScheduler
}

func (FeedScheduler *FeedScheduler) Start() {
	for _, Source := range FeedScheduler.SourceInfoMap {
		FeedScheduler.BaseScheduler.baseSchedulWg.Add(1)

		go FeedScheduler.sourceLoop(Source)
	}

}

func (SourceScheduler *FeedScheduler) sourceLoop(SourceInfo *configs.SourceInfo) {
	defer SourceScheduler.BaseScheduler.baseSchedulWg.Done()

	slog.Info("源调度器已启动：", "source", SourceInfo.SourceName)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			FeedTitleArts, err := scraper.FetchFeed(SourceInfo)
			if err != nil {
				log.Printf("源%s返回文章失败：%v", SourceInfo.SourceName, err)
				continue
			}

			SourceScheduler.FeedTitlePool.Process(SourceInfo, FeedTitleArts)

		case <-SourceScheduler.BaseScheduler.SignalCtx.Done():
			return
		}
	}
}

func (SourceScheduler *FeedScheduler) Close() {
	SourceScheduler.BaseScheduler.Close()
}
