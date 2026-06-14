package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"secmind/internal/model"
	"secmind/internal/scraper"
	"secmind/internal/storage"
	"strings"
	"time"

	"github.com/joho/godotenv" //从文件中获取环境变量用
	"gopkg.in/yaml.v3"         //yaml格式处理
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
}

type ModelSpec struct {
	Name             string   `yaml:"name"`
	APIKeyEnv        string   `yaml:"api_key_env"`
	BaseURLEnv       string   `yaml:"base_url_env"`
	ModelNameEnv     string   `yaml:"model_name_env"`
	SystemPrompt     string   `yaml:"system_prompt"`
	UserPrompt       string   `yaml:"user_prompt"`
	Temperature      float64  `yaml:"temperature"`
	TopP             float64  `yaml:"top_p"`
	MaxTokens        int      `yaml:"max_tokens"`
	FrequencyPenalty float64  `yaml:"frequency_penalty"`
	PresencePenalty  float64  `yaml:"presence_penalty"`
	Stop             []string `yaml:"stop"`
	APIKey           string
	BaseURL          string
	ModelName        string
	ExtraBody        map[string]interface{} `yaml:"extra_body"`
}

func AnalyzeByAI(articles []model.Article, soureceMap map[string]string, articleIndex map[string]*model.Article) {

	//加载环境变量
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Println("加载 .env 失败: ", err)
	}

	//加载yaml配置文件
	var model_param *ModelSpec
	model_param, err = LoadModelConfigByName("configs/model.yaml", "filter")
	if err != nil {
		log.Fatal("加载 .env 失败: ", err)
	}

	//加载提示词
	var promptSys string
	var promptText string
	promptTextdata_sys, err := os.ReadFile(model_param.SystemPrompt)
	if err != nil {
		log.Fatal("读取sys提示词失败: ", err)
	}
	promptTextdata_user, err := os.ReadFile(model_param.UserPrompt)
	if err != nil {
		log.Fatal("读取user提示词失败: ", err)
	}

	promptSys = string(promptTextdata_sys)
	promptText = string(promptTextdata_user)

	for _, a := range articles {
		promptText += fmt.Sprintf("[%s]:%d Title:%s \n", a.Source, a.Id, a.Title)
	}

	var promptMessage []Message
	promptMessage = []Message{
		{Role: "system", Content: promptSys},
		{Role: "user", Content: promptText},
	}

	primaryRespData, err := CallAiApi(model_param, promptMessage)
	if err != nil {
		log.Fatal("Call_Ai返回数据失败: ", err)
	}

	var aiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(primaryRespData, &aiResponse)
	if err != nil {
		log.Fatalf("解析 AI 响应 JSON 失败: %v, 原始响应: %s", err, string(primaryRespData))
	}

	if len(aiResponse.Choices) > 0 {
		fmt.Println("\n=== SecMind 1.0 研判简报 ===")
		fmt.Println(aiResponse.Choices[0].Message.Content)
	}
	
	//从这里开始AI返回的文章信息被解析为ScreenedArticles,因此后续全都应该使用该类型。
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
		fmt.Printf("[%s-%d]：\n%s eng:%s \n%s\n", articleData.Source, articleData.ID, articleData.Title, articleData.EngTitle, articleData.Reason)
		articlehtmldata ,_:= scraper.FetchArticleHtml(realArticle.Link, articleData)

		model_param, err = LoadModelConfigByName("configs/model.yaml", "summarize")
		if err != nil {
			log.Fatal("加载 .env 失败: ", err)
		}

		realArticle.Filename = storage.SaveArticleToMD(articlehtmldata, articleData)

	}

	//临时测试输入选文章功能
	fmt.Print("\n请输入要分析的文章（格式：源缩写-ID，如 TOB-12）：")
	var input string
	_, scanErr := fmt.Scanln(&input)
	if scanErr != nil {
		return
	}

	realArticle, ok := articleIndex[input]
	if !ok {
		fmt.Printf("未找到文章: %s\n", input)
		return
	}

	fmt.Printf("\n已选择文章：%s\n", realArticle.Filename)
	text, err := AnalyzeArticleByAi(model_param, realArticle.Filename)
	if err != nil {
    	fmt.Printf("\n❌ AI分析失败: %v\n", err)
    	return
	}

	//fmt.Printf("\nAI返回的原始内容为：%x\n", text)
	text1 := string(text)
	sourceID := fmt.Sprintf("%s-%d", realArticle.Source, realArticle.Id) 
	storage.SaveArticleToMemory(text1, sourceID)


	err = json.Unmarshal([]byte(text), &aiResponse)
	if err != nil {
		log.Fatalf("解析AI文章分析响应 JSON 失败: %v", err)
	}

	if len(aiResponse.Choices) > 0 {
		fmt.Println("\nAI文章为：")
		fmt.Println(aiResponse.Choices[0].Message.Content)
	}

	//fmt.Println("AI总结摘要内容：", text1)

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

	for i := range wrapper.Models {
		if wrapper.Models[i].Name == modelName {
			// 从环境变量读取真实密钥/URL/模型名
			wrapper.Models[i].APIKey = os.Getenv(wrapper.Models[i].APIKeyEnv)
			wrapper.Models[i].BaseURL = os.Getenv(wrapper.Models[i].BaseURLEnv)
			wrapper.Models[i].ModelName = os.Getenv(wrapper.Models[i].ModelNameEnv)

			return &wrapper.Models[i], nil
		}
	}

	return nil, fmt.Errorf("model %s not found", modelName)
}

func CallAiApi(model_param *ModelSpec, promptMessage []Message) ([]byte, error) {

	reqBody := map[string]interface{}{
		"model":             model_param.ModelName,
		"messages":          promptMessage,
		"temperature":       model_param.Temperature,
		"top_p":             model_param.TopP,
		"max_tokens":        model_param.MaxTokens,
		"frequency_penalty": model_param.FrequencyPenalty,
		"presence_penalty":  model_param.PresencePenalty,
		"stop":              model_param.Stop,
	}

	for k, v := range model_param.ExtraBody {
		reqBody[k] = v
	}

	reqBodyJsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("ai api请求体解析json格式失败", err)
	}

	client := http.Client{Timeout: 60 * time.Second}

	req, err := http.NewRequest("POST", model_param.BaseURL, bytes.NewBuffer([]byte(reqBodyJsonData)))
	if err != nil {
		log.Fatal("创建请求包失败", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+model_param.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errbody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回错误状态 %d: %s", resp.StatusCode, string(errbody))
	}

	primaryRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	//fmt.Println(string(primaryRespBody))
	return primaryRespBody, nil
}
