package analyzer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"secmind/internal/model"
	"secmind/internal/scraper"
	"secmind/internal/storage"
	"strings"

	"github.com/joho/godotenv" //从文件中获取环境变量用
	"gopkg.in/yaml.v3"         //yaml格式处理
)

/*创建 analyzer/types.go：移入 Message, ChatRequest, ModelSpec。

创建 analyzer/client.go：移入 LoadModelConfigByName, CallAiApi，并实现 Client 结构体。

创建 analyzer/filter.go：移入 ParseScreeningTitleJSON，并实现 FilterArticles 方法。

创建 analyzer/summarize.go：移入 AnalyzeArticleByAi 的相关逻辑。

删除或重构 analyzer_title.go：将剩余的通用代码分配到上述文件中。

更新 main 函数：调用新的 Client 接口，并将抓取/保存逻辑移至 article 包。*/


func AnalyzeByAI(articles []model.Article, soureceMap map[string]string, articleIndex map[string]*model.Article) {

	//加载环境变量
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Println("加载 .env 失败1: ", err)
	}

	//加载yaml配置文件
	var model_param *ModelSpec
	model_param, err = LoadModelConfigByName("configs/model.yaml", "filter")
	if err != nil {
		log.Fatal("加载 .env 失败2: ", err)
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


