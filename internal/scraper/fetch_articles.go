package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	readability "codeberg.org/readeck/go-readability/v2"
)

func FetchArticleHtml(articleURL string) string {

	resp, err := http.Get(articleURL)
	if err != nil {
		fmt.Println("请求文章失败，失败原因：", err)
	}
	defer resp.Body.Close()

	pageURL, err := url.Parse(articleURL)
	if err != nil {
		fmt.Println("失败2", err)
	}

	article, err := readability.FromReader(resp.Body, pageURL)
	if err != nil {
		fmt.Println("失败3", err)
	}

	var buf strings.Builder
	article.RenderHTML(&buf)
	//fmt.Printf("%s\n%s\n",article.Title(), buf.String())

	return buf.String()
}
