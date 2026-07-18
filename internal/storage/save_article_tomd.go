package storage

import (
	"fmt"
	"os"
	"secmind/internal/identity"
	"secmind/internal/article"
)

func SaveArticleToMD(htmldata string, articleinfo article.ScreenedArticle) (string){
	var articleTitle string
	articleTitle = "internal/data/articles/" + identity.GenerateFileName(articleinfo)

	articleinfo.ArticleName = articleTitle

	file, err := os.OpenFile(articleTitle, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("文件打开或创建失败：%s\n\n\n", err)
	}
	defer file.Close()

	//fmt.Printf("\n即将存入的文章内容为：%s\n--------------------------------------", htmldata)
	fmt.Fprintf(file, "%s", htmldata)
	
	//暂时用这种笨方法，重构时必须修改此段代码
	var selectedArticle article.Article
	selectedArticle.Filename = articleinfo.ArticleName

	return selectedArticle.Filename
}