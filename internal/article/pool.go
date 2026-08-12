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
)

type SourceManger struct {
	workers map[string]chan FeedArticle
	mu sync.Mutex
}

func TitlePool (ch <-chan FeedArticle) (error) {
	sm := &SourceManger{
		workers: make(map[string]chan FeedArticle),
	}

	for feedArticleTitle := range ch {
		sm.mu.Lock()
		targetCH, exists := sm.workers[feedArticleTitle.Source]
		sm.mu.Unlock()

		feedArticleTitle.ScmFid = fmt.Sprintf("%s-%d", feedArticleTitle.Source, feedArticleTitle.Id)
		feedArticleTitle.FetchedAt = time.Now().UTC()

		if !exists {
			newCH := make(chan FeedArticle, 10)					
			
			sm.mu.Lock()
			sm.workers[feedArticleTitle.Source] = newCH

			go func(source string, newCH chan FeedArticle){
				fetchDate := time.Now().Format("2006-01-02") + ".jsonl"
				inputPath := filepath.Join("internal", "data", "pool", source, fetchDate)
				
				timer := time.NewTimer(3 * time.Minute)
				defer timer.Stop()

				sourcejsonl, err := os.OpenFile(inputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Printf("源%s.jsonl文件打开失败：%v", source, err)
					return 
				}
				defer sourcejsonl.Close()

				for {
					select {
					case art, ok := <-newCH:
						if !ok {
							log.Printf("%s通道不存在", art.Source)
							return
						}
						
						bufioBufferWriter := bufio.NewWriter(sourcejsonl)
						jsonEncoder := json.NewEncoder(bufioBufferWriter)
						err := jsonEncoder.Encode(art)
						if err != nil {
							log.Printf("art按jsonl格式写入jsonEncoder失败\n")
							return 
						}
						
						bufioBufferWriter.Flush()
						sourcejsonl.Sync()

						timer.Reset(3 * time.Minute)

					case <- timer.C:
						sourcejsonl.Sync()
						return
					}
				}
			}(feedArticleTitle.Source, newCH)
			sm.mu.Unlock()

			targetCH = newCH
			targetCH <- feedArticleTitle

		}else {
			targetCH <- feedArticleTitle
		}
	}
	return nil
}