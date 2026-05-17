package storage

import (
	"fmt"
	"os"
	//"secmind/internal/model"
)

func SaveArticleToMD(htmldata string) {
	file, err := os.OpenFile("internal/data/article.html", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("文件打开或创建失败：", err)
	}

	fmt.Fprintf(file, "%s", htmldata)
	
}