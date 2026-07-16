package analyzer

import (
	"fmt"
	"net/http"
	"secmind/configs"
	"strings"
	"time"
)

type Client struct {
	modelSpec  *ModelSpec
	httpClient *http.Client
}

func (c *Client) Execute(templatestName string, userInput string)(string, error){
	
}

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

	return client, fmt.Errorf("modelSpec核心字段为空：%w", err)
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
