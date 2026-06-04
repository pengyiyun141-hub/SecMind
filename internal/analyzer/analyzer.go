package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"secmind/internal/model"
	"secmind/internal/scraper"
	"secmind/internal/storage"
	"strings"

	"github.com/joho/godotenv" //从文件中获取环境变量用
	"gopkg.in/yaml.v3"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ModelSpec struct {
    Name            string    `yaml:"name"`
    APIKeyEnv       string    `yaml:"api_key_env"`
    BaseURLEnv      string    `yaml:"base_url_env"`
    ModelNameEnv    string    `yaml:"model_name_env"`
    SystemPrompt    string    `yaml:"system_prompt"`
    UserPrompt      string    `yaml:"user_prompt"`
    Temperature     float64   `yaml:"temperature"`
    TopP            float64   `yaml:"top_p"`
    MaxTokens       int       `yaml:"max_tokens"`
    FrequencyPenalty float64  `yaml:"frequency_penalty"`
    PresencePenalty float64   `yaml:"presence_penalty"`
    Stop            []string  `yaml:"stop"`
    // 运行时填充的真实值（不从 YAML 读）
    APIKey          string
    BaseURL         string
    ModelName       string
}

func AnalyzeByAI(articles []model.Article, soureceMap map[string]string, articleIndex map[string]*model.Article) {

	err := godotenv.Load("configs/.env")

	if err != nil {
		log.Println("加载 .env 失败: ", err)
	}

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	apiURL := os.Getenv("DEEPSEEK_API_URL")
	modelName := os.Getenv("DEEPSEEK_MODEL")

	var promptText string
	for _, a := range articles {
		promptText += fmt.Sprintf("id:[%d] Title:%s Source:%s\n", a.Id, a.Title, a.Source)
	}

	requestPayload := ChatRequest{
		Model: modelName,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一名资深的网络安全的研究员，请对这份列表中的标题进行筛选，选出你认为跟AI最相关的五篇，并说明理由。",
			},
			{
				Role:    "user",
				Content: "请将原标题翻译成中文，要求保留原英文标题。要求返回结果为json格式，其中包含字段：ID,Title,engtitle（英文标题）,Source,Reason。其中id只显示数字，且这里的id是我传进来的字段，另外id字段是int类型，不要带引号，返回结果时请不要自行生成ID。source返回我传入的缩写即可。这是今天的论文列表：\n" + promptText,
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

	screeneddata, err := ParseScreeningTitleJSON(aiResponse.Choices[0].Message.Content)
	if err != nil {
		fmt.Println("第一轮内容筛选解析失败", err)
	}

	for _, articleData := range screeneddata {

		key := fmt.Sprintf("%s-%d", articleData.Source, articleData.ID)
		realArticle, ok := articleIndex[key]
		if !ok {
			fmt.Printf("未找到文章: %s\n", key)
			continue
		}
		fmt.Println("开始抓取文章：", realArticle.Link)
		fmt.Printf("[%s-%d]：\n%s eng:%s \n%s\n\n", articleData.Source, articleData.ID, articleData.Title, articleData.EngTitle, articleData.Reason)
		articlehtmldata := scraper.FetchArticleHtml(realArticle.Link)

		storage.SaveArticleToMD(articlehtmldata, articleData.EngTitle)
	}

}

func extractJSON(text string) string {

	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")

	if start == -1 || end == -1 || start >= end {
		fmt.Println("start,end分别为：", start, end)
		return ""
	}

	return text[start : end+1]
}

func ParseScreeningTitleJSON(content string) ([]model.ScreenedArticle, error) {

	jsonStr := extractJSON(content)
	if jsonStr == "" {
		return nil, fmt.Errorf("AI回复中未找到有效的json数组")
	}

	var screeneddata []model.ScreenedArticle
	err := json.Unmarshal([]byte(jsonStr), &screeneddata)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	return screeneddata, err
}

func LoadModelConfigByName(yamlPath, modelName string) (*ModelSpec, error) {
	yamlfile, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("加载model.yaml文件失败: %w", err)
	}

	var wrapper struct {
		Models []ModelSpec `yaml:"models"`
	}

	err = yaml.Unmarshal(yamlfile, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("解析yaml值失败: %w", err)
	}

	for _, m := range wrapper.Models {
		if m.Name == modelName {
			// 从环境变量读取真实密钥/URL/模型名
            m.APIKey = os.Getenv(m.APIKeyEnv)
            m.BaseURL = os.Getenv(m.BaseURLEnv)
            m.ModelName = os.Getenv(m.ModelNameEnv)
            return &m, nil
		}
	}
	
	return nil, fmt.Errorf("model %s not found", modelName)
}