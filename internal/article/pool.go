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

type Pool struct {
	wg 		sync.WaitGroup
	poolCfg	*configs.PoolConfigs
}

type SourceManager struct {
	workers map[string]chan FeedArticle
	mu      sync.Mutex
}

func NewPool(poolCfg *configs.PoolConfigs) (*Pool, error) {
	pool := &Pool{
		poolCfg: poolCfg,
		wg: sync.WaitGroup{},
	}
	
	return pool, nil
}

var wg sync.WaitGroup

func TitlePool(workerCH <-chan FeedArticle) error {
	sm := &SourceManager{
		workers: make(map[string]chan FeedArticle),
	}

	for feedArticleTitle := range workerCH {

		feedArticleTitle.ScmFid = fmt.Sprintf("%s-%d", feedArticleTitle.Source, feedArticleTitle.Id)
		feedArticleTitle.FetchedAt = time.Now().UTC()

		target := sm.getOrCreateWorker(feedArticleTitle.Source)
		target <- feedArticleTitle
	}

	wg.Wait()
	return nil
}

func (sm *SourceManager) getOrCreateWorker(source string) chan FeedArticle {

	sm.mu.Lock()
	targetCH, exists := sm.workers[source]

	if !exists {
		newCH := make(chan FeedArticle, 10)
		sm.workers[source] = newCH
		targetCH = newCH
		wg.Add(1)
	}
	sm.mu.Unlock()

	if !exists {
		go func(source string, newCH chan FeedArticle) {
			fetchDate := time.Now().Format("2006-01-02") + ".jsonl"
			inputPath := filepath.Join("internal", "data", "pool", source, fetchDate)

			timer := time.NewTimer(3 * time.Minute)
			defer timer.Stop()

			err := os.MkdirAll(filepath.Dir(inputPath), 0755)
			if err != nil {
				log.Printf("创建目录失败：%v", err)
				return
			}

			sourcejsonl, err := os.OpenFile(inputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Printf("源%s.jsonl文件打开失败：%v", source, err)
				return
			}
			defer sourcejsonl.Close()

			bufioBufferWriter := bufio.NewWriter(sourcejsonl)
			jsonEncoder := json.NewEncoder(bufioBufferWriter)

			for {
				select {
				case art, ok := <-newCH:
					if !ok {
						log.Printf("%s通道不存在", art.Source)
						return
					}

					err := jsonEncoder.Encode(art)
					if err != nil {
						log.Printf("art按jsonl格式写入jsonEncoder失败\n")
						return
					}

					bufioBufferWriter.Flush()
					timer.Reset(3 * time.Minute)

				case <-timer.C:
					sourcejsonl.Sync()
					wg.Done()
					return
				}
			}
		}(source, targetCH)

		return targetCH
	} else {
		return targetCH
	}
}
