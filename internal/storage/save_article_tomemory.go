package storage

import (
	"fmt"
	"os"
)

func SaveArticleToMemory(htmldata string, sourceID string) error {
	var articleTitlePath string
	articleTitlePath = "internal/data/memory/" + sourceID

	file, err := os.OpenFile(articleTitlePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("文件打开或创建失败：", err)
	}

	fmt.Fprintf(file, "%s", htmldata)

	return err
}
