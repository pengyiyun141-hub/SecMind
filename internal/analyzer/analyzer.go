package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"secmind/internal/model"
	"github.com/joho/godotenv" //从文件中获取环境变量用
)

type Article_AI struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func analyzeByAI(articles []model.Article) {

	err := godotenv.Load("../../configs/.env")

	if err != nil {
		log.Println("加载 .env 失败: ", err)
	}

	apiKey := os.Getenv("ZHIPU_API_KEY")
	apiURL := os.Getenv("ZHIPU_API_URL")
	modelName := os.Getenv("ZHIPU_MODEL")

	var ai_message_result []Article_AI
	for _, ai_message := range articles {

		ai_message := Article_AI{
			Id:     ai_message.Id,
			Title:  ai_message.Title,
			Link:   ai_message.Link,
			Source: ai_message.Source,
		}

		ai_message_result = append(ai_message_result, ai_message)
	}

	var promptText string

	for _, a := range articles {
		promptText += fmt.Sprintf("id:[%d] Title:%s Source:%s\n", a.Id, a.Title, a.Source)
	}

	requestPayload := ChatRequest{
		Model: modelName,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一名资深的网络安全的研究员，请对这份列表中的标题进行筛选，选出你认为与智能安全最相关的十篇，并说明理由。",
			},
			{
				Role:    "user",
				Content: "请将原标题翻译成中文，并标注source和id，要求保留原标题。这是今天的论文列表：\n" + promptText,
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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

	if len(aiResponse.Choices) > 0 {
		fmt.Println("\n=== SecMind 1.0 研判简报 ===")
		fmt.Println(aiResponse.Choices[0].Message.Content)
	}

}
