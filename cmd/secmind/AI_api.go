package main

import(
	"fmt"
	"log"
	"json/encoding"
)

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
