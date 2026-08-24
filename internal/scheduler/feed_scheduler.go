package scheduler

import (
	"log"
	"secmind/configs"
	"secmind/internal/article"
	"secmind/internal/scraper"
	"time"
)

type FeedScheduler struct {
	SourceInfoMap	map[string]*configs.SourceInfo
	FeedTitlePool	*article.FeedTitlePool
	BaseScheduler	BaseScheduler
}

func NewFeedScheduler(SourceInfoMap map[string]*configs.SourceInfo, FeedTitlePool *article.FeedTitlePool) (*FeedScheduler) {
	feedScheduler := &FeedScheduler{
		SourceInfoMap: SourceInfoMap,
		FeedTitlePool: FeedTitlePool,
		BaseScheduler: *NewBaseScheduler(),
	}

	return feedScheduler
}

func (FeedScheduler *FeedScheduler)Start() () {
	for _, Source := range FeedScheduler.SourceInfoMap {
		FeedScheduler.BaseScheduler.baseSchedulWg.Add(1)

		go FeedScheduler.sourceLoop(Source)
	}

}

func (SourceScheduler *FeedScheduler)sourceLoop(SourceInfo *configs.SourceInfo) () {
	defer SourceScheduler.BaseScheduler.baseSchedulWg.Done()
	for{
		select{
		case <-time.After(1800):
			FeedTitleArts, err := scraper.FetchFeed(SourceInfo)
			if err != nil {
				log.Printf("源%s返回文章失败：%v", SourceInfo.SourceName, err)
				continue
			}

			SourceScheduler.FeedTitlePool.Process(SourceInfo, FeedTitleArts)

		case <-SourceScheduler.BaseScheduler.ctx.Done():
			return 
		}
	}
}