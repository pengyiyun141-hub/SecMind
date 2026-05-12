package main

import(
	"os"
	"fmt"
)

func SaveToMD(articles []Article) (error){
	file, err := os.OpenFile("../../internal/storage/intel_report.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("文件打开或创建失败：",err)
	}
	
	for _, text := range articles{
		fmt.Fprintf(file, "### [%s]\n(%s)\n\n", text.Title, text.Link)
	}
	
	return err
}
