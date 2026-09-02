package article

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"secmind/configs"
	"sync"
	"time"
	"log/slog"
)

//新架构

type FeedTitlePool struct {
	feedtitlePoolWg sync.WaitGroup
	poolCfg         *configs.PoolConfigs
	lastSuccessAt	map[string]time.Time
	mu      		sync.Mutex
}

type FeedTitleJob struct {
	TitlePoolCfg *configs.PoolConfigs
	SourceInfo   *configs.SourceInfo
	FeedArticles []FeedArticle
}

func NewPool(poolCfg *configs.PoolConfigs) (*FeedTitlePool, error) {
	TitlePool := &FeedTitlePool{
		poolCfg:         poolCfg,
		feedtitlePoolWg: sync.WaitGroup{},
		lastSuccessAt: 	 make(map[string]time.Time),
	}

	return TitlePool, nil
}

func (FeedTitlePool *FeedTitlePool) Process(SourceInfo *configs.SourceInfo, FeedArticles []FeedArticle) (int, error) {
	FeedTitlePool.feedtitlePoolWg.Add(1)

	fileDedup, _ := NewFileDedup(SourceInfo.SourceName)
	dedupFeedArticles := fileDedup.Filter(FeedArticles)

	FeedTitleJob := &FeedTitleJob{
		TitlePoolCfg: FeedTitlePool.poolCfg,
		SourceInfo:   SourceInfo,
		FeedArticles: dedupFeedArticles,
	}

	go func() {
		defer FeedTitlePool.feedtitlePoolWg.Done()
		err := FeedTitlePool.Persist(FeedTitleJob)
		if err != nil {
			slog.Error("[process()]:Feed源批次持久化失败", "source", SourceInfo.SourceName, "error", err) 
		}
	}()

	return len(FeedArticles), nil
}

func (FeedTitlePool *FeedTitlePool) Persist(FeedTitleJob *FeedTitleJob) error {
	//先判断是否有新文章
	if len(FeedTitleJob.FeedArticles) == 0 {
		slog.Info("[Process()]:Feed源未更新，", "source", FeedTitleJob.SourceInfo.SourceName)
		return nil
	}

	FeedTitlePool.mu.Lock()
	now := time.Now()
	last, ok := FeedTitlePool.lastSuccessAt[FeedTitleJob.SourceInfo.SourceName]
	var interval string
	if ok && !last.IsZero() {
    	interval = now.Sub(last).String()
	}else {
		interval = "first"
	}
	FeedTitlePool.lastSuccessAt[FeedTitleJob.SourceInfo.SourceName] = now
	FeedTitlePool.mu.Unlock()

	fetchDate := time.Now().Format("2006-01-02") + ".jsonl"
	inputPath := filepath.Join(FeedTitlePool.poolCfg.DataDir, FeedTitleJob.SourceInfo.SourceName, fetchDate)

	err := os.MkdirAll(filepath.Dir(inputPath), 0755)
	if err != nil {
		return fmt.Errorf("[Persist()]:创建目录失败 %s:%w", inputPath, err)
	}

	sourcejsonl, err := os.OpenFile(inputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("[Persist()]:源%s.jsonl文件打开失败：%w", FeedTitleJob.SourceInfo.SourceName, err)
	}
	defer sourcejsonl.Close()

	bufioBufferWriter := bufio.NewWriter(sourcejsonl)
	jsonEncoder := json.NewEncoder(bufioBufferWriter)

	new_count := 0
	for _, feedArticleTitle := range FeedTitleJob.FeedArticles {
		feedArticleTitle.ScmFid = fmt.Sprintf("%s-%d", feedArticleTitle.Source, feedArticleTitle.Id)
		feedArticleTitle.FetchedAt = time.Now().UTC()

		err := jsonEncoder.Encode(feedArticleTitle)
		if err != nil {
			slog.Error("[Persist()]:条目写入失败\n", "source", feedArticleTitle.Source, "Link", feedArticleTitle.Link, "error", err)
			continue
		}

		new_count++
		bufioBufferWriter.Flush()
	}
	sourcejsonl.Sync()
	slog.Info("[Persist()]:", "source", FeedTitleJob.SourceInfo.SourceName, "new_count:", new_count, "interval:", interval)

	return nil
}

func (TitlePool *FeedTitlePool) Close() {
	TitlePool.feedtitlePoolWg.Wait()
}
