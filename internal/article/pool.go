package article

import (
	"fmt"
	"path/filepath"
	"time"
)

func TitlePool (ch <-chan FeedArticle) (error) {
	var inputPath string
	fetchData := time.Now().Format("2006-01-02")

	for feedarticle := range ch {
		inputPath = filepath.Join("data", "pool", feedarticle.Source, fetchData)
		fmt.Fprintf(inputPath,,)  //遍历ch中的结构体写进对应的文件
		//tomic.AddInt32(&counter, 1) 使用原子操作为特定源的ID计数
	}
	
		
}