package storage

import (
	"fmt"
	"os"
	//"secmind/internal/model"
)

func SaveArticleToMD(htmldata string, title string) {
	var articleTitlePath string
	articleTitlePath = "internal/data/articles/" + title

	file, err := os.OpenFile(articleTitlePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("文件打开或创建失败：", err)
	}

	fmt.Fprintf(file, "%s", htmldata)
	
}