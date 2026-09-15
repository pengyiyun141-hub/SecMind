package configs

import (
	
)

type SecmindConfigs struct {
	Aiconfigs   *AiConfigs
	Feedconfigs *FeedConfigs
	Poolconfigs *PoolConfigs
	Logconfigs	*LogConfigs
}

// AI配置
type AiConfigs struct {
	Apiinfo    map[string]*ApiInfo
	Promptinfo map[string]*PromptInfo
	Modelinfo  map[string]*ModelInfo
}

type ApiInfo struct { //注意该结构体暂时没用到，Api的相关信息直接被填入model中了。
	Baseurl   string
	Modelname string
	Apikey    string
}

type PromptInfo struct {
	System string
	User   string
}

type ModelInfo struct {
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
	PromptSysText    string
	PromptUsrText    string
	ExtraBody        map[string]interface{} `yaml:"extra_body"`
}

// Feed配置
type FeedConfigs struct {
	SourceInfoMap map[string]*SourceInfo
}

type SourceInfo struct {
	SourceName string `json:"SourceName"`
	Kind	   string `json:"Kind"`
	Type       string `json:"Type"`
	Schedule   string `json:"Schedule"`
	Enabled    bool   `json:"Enabled"`

	Feed *FeedSourceInfo `json:"feed,omitempty"`
    API  *APISourceInfo  `json:"api,omitempty"`
}

type FeedSourceInfo struct {
    URL  		string `json:"url"`
    Type 		string `json:"type"` // rss / atom
}

type APISourceInfo struct {
    Provider string `json:"provider"`
    BaseURL  string `json:"base_url"`
    Search   string `json:"search"`
    Sort     string `json:"sort"`
    PerPage  int    `json:"per_page"`
    Select   string `json:"select"`
	Apikey	 string `json:"api_key"`
}

type PoolConfigs struct {
	DataDir string `json:"DataDir"`
}

type LogConfigs struct {
	Level  string		`yaml:"level"`
    Format string   	`yaml:"format"`
    Output string   	`yaml:"output"`
    File   string   	`yaml:"file"`
}