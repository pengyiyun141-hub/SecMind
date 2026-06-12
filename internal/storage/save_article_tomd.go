package storage

import (
	"fmt"
	"os"
	//"secmind/internal/model"
)

func SaveArticleToMD(htmldata string, title string) {
	var articleTitlePath string
	articleTitlePath = "internal/data/articles/" + title

	file, err := os.OpenFile(articleTitlePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("文件打开或创建失败：%s\n\n\n", err)
	}
	defer file.Close()

	fmt.Printf("\n即将存入的文章内容为：%s\n--------------------------------------", htmldata)
	fmt.Fprintf(file, "%s", htmldata)
	
}