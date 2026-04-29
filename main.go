package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/PuerkitoBio/goquery"
	//"golang.org/x/text/message"
	"bytes"
	"encoding/json"
	//"io"
)

type Article struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

type Message struct {
	Role    string `json:"role"`    
	Content string `json:"content"` 
}

type ChatRequest struct {
	Model    string    `json:"model"`    
	Messages []Message `json:"messages"` 
}

func analyzeByAI(articles []Article){
	
	apiKey := "" 
	apiURL := ""

	var promptText string

	for _,a := range articles{
		promptText += fmt.Sprintf("[%d] %s\n",a.ID,a.Title)
	}

	fmt.Println(promptText)

	requestPayload := ChatRequest{
		Model : "glm-4-flash",
		Messages : []Message{
			{
				Role:    "system",
				Content: "你是一个说话极其刻薄的资深黑客，请用讽刺的口吻评价这些论文的实战价值",
			},
			{
				Role:    "user", 
				Content: "请将标题翻译成中文，这是今天的论文列表：\n" + promptText,
			},
		},
	}

		jsonData, err := json.Marshal(requestPayload)

		if err != nil {
			fmt.Println("打包失败:", err)
			return
		}

		//fmt.Println("\n=== 准备发送给 AI 的 JSON 包预览 ===")
		//fmt.Println(string(jsonData))

		bodyPipe := bytes.NewBuffer(jsonData)

		req, err := http.NewRequest("POST", apiURL, bodyPipe)
	if err != nil {
		log.Fatal("构造请求对象失败:", err)
	}

	// 6. 贴邮票：设置 Header（智谱/OpenAI 接口必须要求这两项）
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey) // Bearer 后面有一个空格

	client := &http.Client{}
	fmt.Println("\n[AI] 正在进行情报分析，请稍候...")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("向 AI 发送请求失败:", err)
	}
	defer resp.Body.Close()

	var aiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}	

	err = json.NewDecoder(resp.Body).Decode(&aiResponse)
	if err != nil {
		log.Fatal("解析 AI 回复失败:", err)
	}

	// 9. 最终战果展示
	if len(aiResponse.Choices) > 0 {
		fmt.Println("\n=== SecMind 1.0 研判简报 ===")
		fmt.Println(aiResponse.Choices[0].Message.Content)
	}
}

	
	

func main() {
	// arXiv 的计算机安全最近更新列表
	url := "https://arxiv.org/list/cs.CR/recent"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("状态码错误: %d", resp.StatusCode)
	}

	
	// 2. 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	htmlContent, _ := doc.Html()
	fmt.Println("还原后的树结构预览：\n", htmlContent[:500])

	var reportList []Article

	// 3. arXiv 的结构：标题在 <div class="list-title"> 里面
	fmt.Println("=== 正在捕获 arXiv 全球安全研究前沿 ===")
	
	doc.Find("dt").Each(func(i int, dt *goquery.Selection) {
		if i >= 10 { return } // 我们先抓前 10 篇

		// 寻找对应的标题 (在下一个兄弟节点 dd 中)
		dd := dt.Next()
		title := dd.Find(".list-title").Text()
		// 这里先去掉title两边的
		title = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(title), "Title:"))

		// 寻找 PDF 链接 (在当前 dt 节点中)
		pdfPath, _ := dt.Find("a[title='Download PDF']").Attr("href")
		pdfLink := "https://arxiv.org" + pdfPath

		item := Article{
			ID:    i + 1,
			Title: title,
			Link:  pdfLink,
		}
		reportList = append(reportList, item)
		
		fmt.Printf("[%d] 发现论文: %s\n", item.ID, item.Title)
		fmt.Printf("PDF链接: %s\n\n",  item.Link)
	})

	fmt.Printf("\n--- 成功！共抓取 %d 篇高价值论文数据 ---\n", len(reportList))

	analyzeByAI(reportList)
}