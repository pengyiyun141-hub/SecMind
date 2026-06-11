package storage

import (
	"fmt"
	"os"
	//"secmind/internal/model"
)

func SaveArticleToMemory(htmldata string, link string) error {
	var articleTitlePath string
	articleTitlePath = "internal/data/memory/" + link

	file, err := os.OpenFile(articleTitlePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("文件打开或创建失败：", err)
	}

	fmt.Fprintf(file, "%s", htmldata)

	return err
}
