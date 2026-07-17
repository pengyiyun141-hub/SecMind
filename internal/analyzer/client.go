package analyzer

import (
	"fmt"
	"net/http"
	"secmind/configs"
	"strings"
	"time"
	"encoding/json"
	"io"
	"bytes"
)

/*
func (client *Client) Execute(userInput string)([]byte, error){
	client.modelSpec.PromptUserText = client.modelSpec.PromptUserText + userInput
	 primitivedata, err := CallAiApi(client)
	 if err != nil {
		return nil, fmt.Errorf("Execute：%w", err)
	 }

	 return primitivedata, err 
}*/

func NewClient(role string, AiCfgs *configs.AiConfigs)(*Client, error){
	rolesplit := strings.SplitN(role, "-", 2)
	modelCfg := AiCfgs.Modelinfo[rolesplit[0]]
	modelCfg.SystemPrompt = AiCfgs.Promptinfo[role].System
	modelCfg.PromptSysText = AiCfgs.Promptinfo[role].User

	var ms ModelSpec
	ms = ModelSpec{
		role: role,
		Temperature: AiCfgs.Modelinfo[role].Temperature,
		TopP: AiCfgs.Modelinfo[role].TopP,
		MaxTokens: AiCfgs.Modelinfo[role].MaxTokens,
		FrequencyPenalty: AiCfgs.Modelinfo[role].FrequencyPenalty,
		PresencePenalty: AiCfgs.Modelinfo[role].PresencePenalty,
		Stop: AiCfgs.Modelinfo[role].Stop,
		APIKey: AiCfgs.Modelinfo[role].APIKey,
		BaseURL: AiCfgs.Modelinfo[role].BaseURL,
		ModelName: AiCfgs.Modelinfo[role].ModelName,
		PromptSystemText: modelCfg.PromptSysText,
		PromptUserText: modelCfg.PromptUsrText,
		ExtraBody: AiCfgs.Modelinfo[role].ExtraBody,
	}

	var hc http.Client
	hc = http.Client{
		Timeout: 30 * time.Second,  // 整个请求（含连接、发送、接收）的总超时
    	Transport: &http.Transport{
        	MaxIdleConns:        100,              // 最大空闲连接数（所有 host 合计）
        	MaxIdleConnsPerHost: 10,               // 每个 host 的最大空闲连接数
        	IdleConnTimeout:     90 * time.Second, // 空闲连接存活时间
        	TLSHandshakeTimeout: 10 * time.Second, // TLS 握手超时
        	ExpectContinueTimeout: 1 * time.Second,
		},
	}

	client := &Client {
		modelSpec: &ms,
		httpClient: &hc,
	}

	err := client.Validate()
	if err != nil {
		return nil, fmt.Errorf("modelSpec核心字段为空：%w", err)
	}

	return client, err
}

func (client *Client)CallAiApi(userInput string) ([]byte, error) {
	userPrompt := client.modelSpec.PromptUserText + userInput
	chatrequest := &ChatRequest {
		Model: client.modelSpec.ModelName,
		Message: []Message{
			{Role: "system", Content: client.modelSpec.PromptSystemText},
			{Role: "user", Content: userPrompt},
		},
		Temperature: client.modelSpec.Temperature,
		TopP: client.modelSpec.TopP,
		MaxTokens: client.modelSpec.MaxTokens,
		FrequencyPenalty: client.modelSpec.FrequencyPenalty,
		PresencePenalty: client.modelSpec.PresencePenalty,
		Stop: client.modelSpec.Stop,
	}

	chatrequestJsonData, err := json.Marshal(chatrequest)
	if err != nil {
		return nil, fmt.Errorf("CallAiApi请求体解析json格式失败", err)
	}

	reqclient := client.httpClient

	req, err := http.NewRequest("POST", client.modelSpec.BaseURL, bytes.NewBuffer([]byte(chatrequestJsonData)))
	if err != nil {
		return nil, fmt.Errorf("NewRequest创建请求包失败：%w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ client.modelSpec.APIKey)

	resp, err := reqclient.Do(req)
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

func (c *Client) Validate() error {
    if c.modelSpec.role == "" {
        return fmt.Errorf("role不能为空;")
    }
    if c.modelSpec.APIKey == "" {
        return fmt.Errorf("APIKey不能为空;")
    }
    if c.modelSpec.BaseURL == "" {
        return fmt.Errorf("BaseURL不能为空;")
    }
    if c.modelSpec.PromptSystemText == "" {
        return fmt.Errorf("PromptSystemText不能为空;")
    }
    if c.modelSpec.PromptUserText == "" {
        return fmt.Errorf("PromptUserText不能为空;")
    }

    return nil
}
