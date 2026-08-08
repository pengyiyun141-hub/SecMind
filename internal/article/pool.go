package article

import (
	"encoding/json"
	"fmt"
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
	fetchData := time.Now().Format("2006-01-02")
	sm := &SourceManger{
		workers: make(map[string]chan FeedArticle),
	}

	for feedArticleTitle := range ch {
		sm.mu.Lock()
		sCH, exists := sm.workers[feedArticleTitle.Source]
		sm.mu.Unlock()

		if !exists {
			newCH := make(chan FeedArticle, 10)
			go func(source string, newCH chan FeedArticle) {
				inputPath := filepath.Join("data", "pool", feedArticleTitle.Source, fetchData, ".jsonl")
				
				for art := range newCH {
					fmt.Fprintf(inputPath, "%s%d-%s:%s", feedArticleTitle.Source, feedArticleTitle.Id, feedArticleTitle.Title, feedArticleTitle.Link)
				}
			}(feedArticleTitle.Source, newCH)
		}
	}
			
}