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
}

type FeedTitleJob struct {
	TitlePoolCfg *configs.PoolConfigs
	SourceInfo   *configs.SourceInfo
	FeedArticles []FeedArticle
}

type SourceManager struct {
	workers map[string]chan FeedArticle
	mu      sync.Mutex
}

func NewPool(poolCfg *configs.PoolConfigs) (*FeedTitlePool, error) {
	TitlePool := &FeedTitlePool{
		poolCfg:         poolCfg,
		feedtitlePoolWg: sync.WaitGroup{},
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

func (TitlePool *FeedTitlePool) Persist(FeedTitleJob *FeedTitleJob) error {
	//先判断是否有新文章
	if len(FeedTitleJob.FeedArticles) == 0 {
		slog.Info("[Process()]:Feed源未更新，", "source", FeedTitleJob.SourceInfo.SourceName)
		return nil
	}

	fetchDate := time.Now().Format("2006-01-02") + ".jsonl"
	inputPath := filepath.Join(TitlePool.poolCfg.DataDir, FeedTitleJob.SourceInfo.SourceName, fetchDate)

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

	count := 0
	for _, feedArticleTitle := range FeedTitleJob.FeedArticles {
		feedArticleTitle.ScmFid = fmt.Sprintf("%s-%d", feedArticleTitle.Source, feedArticleTitle.Id)
		feedArticleTitle.FetchedAt = time.Now().UTC()

		err := jsonEncoder.Encode(feedArticleTitle)
		if err != nil {
			slog.Error("[Persist()]:条目写入失败\n", "source", feedArticleTitle.Source, "Link", feedArticleTitle.Link, "error", err)
			continue
		}

		count++
		bufioBufferWriter.Flush()
	}
	sourcejsonl.Sync()
	slog.Info("[Persist()]:", "source", FeedTitleJob.SourceInfo.SourceName, "本次抓取共持久化文章数量:", count)

	return nil
}

func (TitlePool *FeedTitlePool) Close() {
	TitlePool.feedtitlePoolWg.Wait()
}
