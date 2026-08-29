package article

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"secmind/configs"
)
//新架构

type FeedTitlePool struct {
	feedtitlePoolWg 	sync.WaitGroup
	poolCfg		*configs.PoolConfigs
}

type FeedTitleJob struct {
	TitlePoolCfg	*configs.PoolConfigs
	SourceInfo		*configs.SourceInfo
	FeedArticles	[]FeedArticle
}

type SourceManager struct {
	workers map[string]chan FeedArticle
	mu      sync.Mutex
}

func NewPool(poolCfg *configs.PoolConfigs) (*FeedTitlePool, error) {
	TitlePool := &FeedTitlePool{
		poolCfg: poolCfg,
		feedtitlePoolWg: sync.WaitGroup{},
	}
	
	return TitlePool, nil
}

func (FeedTitlePool *FeedTitlePool) Process(SourceInfo *configs.SourceInfo, FeedArticles []FeedArticle) (int, error) {
	FeedTitlePool.feedtitlePoolWg.Add(1)

	FeedTitleJob := &FeedTitleJob{
		TitlePoolCfg: FeedTitlePool.poolCfg,
		SourceInfo: SourceInfo,
		FeedArticles: FeedArticles,
	}

	go func() {
		defer FeedTitlePool.feedtitlePoolWg.Done()
		err := FeedTitlePool.Persist(FeedTitleJob)
		if err != nil {
			log.Printf("")		//暂时没想好写什么
		}
	}()

	return len(FeedArticles), nil
}  

func (TitlePool *FeedTitlePool) Persist(FeedTitleJob *FeedTitleJob) (error) {
	fetchDate := time.Now().Format("2006-01-02") + ".jsonl"
	inputPath := filepath.Join(TitlePool.poolCfg.DataDir, FeedTitleJob.SourceInfo.SourceName, fetchDate)

	err := os.MkdirAll(filepath.Dir(inputPath), 0755)
	if err != nil {
		log.Printf("创建目录失败：%v", err)
		return nil
	}

	sourcejsonl, err := os.OpenFile(inputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("源%s.jsonl文件打开失败：%v", FeedTitleJob.SourceInfo.SourceName, err)
		return nil
	}
	defer sourcejsonl.Close()
	
	bufioBufferWriter := bufio.NewWriter(sourcejsonl)
	jsonEncoder := json.NewEncoder(bufioBufferWriter)

	for _, feedArticleTitle := range FeedTitleJob.FeedArticles {
		feedArticleTitle.ScmFid = fmt.Sprintf("%s-%d", feedArticleTitle.Source, feedArticleTitle.Id)
		feedArticleTitle.FetchedAt = time.Now().UTC()

		err := jsonEncoder.Encode(feedArticleTitle)
		if err != nil {
			log.Printf("[err]%s-%d条目写入失败\n", feedArticleTitle.Source, feedArticleTitle.Id)
		}

		bufioBufferWriter.Flush()
	}
	sourcejsonl.Sync()

	return nil
}

func (TitlePool *FeedTitlePool) Close() {
	TitlePool.feedtitlePoolWg.Wait()
}