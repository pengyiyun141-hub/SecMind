package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	//"text/scanner"
	//"strings"
	//"github.com/PuerkitoBio/goquery"
	//"golang.org/x/text/message"
	"bytes"
	"encoding/json"
	//"io"
	"os"
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
				Content: "你是一名资深的网络安全工程师，请对这份列表中的标题进行筛选，选出你认为的最前沿，最理论的，并给出判断。",
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
	// 读取urls.txt
	urls_file, err:= os.Open("../../configs/urls.txt")

	if err != nil {
		log.Fatal("文件打开失败:", err)
	}

	scanner := bufio.NewScanner(urls_file)

	var urls_str []string
	for scanner.Scan() {
		line := scanner.Text()

		urls_str = append(urls_str, line)
	}
	
		var xmlData_slice []Article1

		for article := range Fetch(urls_str) {
    		xmlData_slice = append(xmlData_slice, article)
		}

	i := len(xmlData_slice)

	if  i > 0 {
		

		for _, article := range xmlData_slice {
    		fmt.Printf("标题 %d: %s\n[%s] 源:[%s]\n\n", article.Id+1, article.Title, article.Link, article.Source)
		}

		for a, article := range xmlData_slice {
    		fmt.Printf("标题 %d: %s\n[%s]\n\n", a+1, article.Title, article.Link)

		}
	}
	
	fmt.Println("")
	
	SaveToMD(xmlData_slice)

}	