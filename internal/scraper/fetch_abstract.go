package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

func FetchAbstract(articleURL string) (string, error) {
	
	resp, err := http.Get(articleURL)
	if err != nil {
		return "", fmt.Errorf("fetchAbstract功能请求%s文章摘要失败。", articleURL)
	}
	defer resp.Body.Close()

	articleAbstractDoc, err:= goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("goquery.NewDocumentFromReader功能解析html失败。")
	}

	
	articleAbstract := strings.TrimSpace(articleAbstractDoc.Find("blockquote.abstract").First().Text())
	//fmt.Printf("%s的articleAbstract:\n%s\n------------------------------------------------\n\n", articleURL, articleAbstract)

	return articleAbstract, err
}
