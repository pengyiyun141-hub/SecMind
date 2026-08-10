package article

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"log"
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

		if !exists {
			newCH := make(chan FeedArticle, 10)

			go func(source string, newCH chan FeedArticle){
				fetchData := time.Now().Format("2006-01-02") + ".jsonl"
				inputPath := filepath.Join("data", "pool", source, fetchData)

				sourcejsonl, err := os.OpenFile(inputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					log.Printf("源%s.jsonl文件打开失败：%v", source, err)
					return 
				}
				defer sourcejsonl.Close()

				for art := range newCH {
					fmt.Fprintf(sourcejsonl, "%s%d-%s:%s", art.Source, art.Id, art.Title, art.Link)
				}

				return
			}(feedArticleTitle.Source, newCH)

			sm.mu.Lock()
			sm.workers[feedArticleTitle.Source] = newCH
			sm.mu.Unlock()

			newCH <- feedArticleTitle

		}else {
			targetCH <- feedArticleTitle
		}
	}
	return nil
}