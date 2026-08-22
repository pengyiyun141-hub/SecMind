package scheduler

import (
	"secmind/configs"
	"secmind/internal/article"
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

}