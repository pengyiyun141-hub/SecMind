package configs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

//所有配置
type SecmindConfigs struct{
	aiconfigs		*AiConfigs
	feedconfigs     *FeedConfigs 
}

//AI配置
type AiConfigs struct{
	apiinfo		map[string]*ApiInfo
	promptinfo	map[string]*PromptInfo
	modelinfo	map[string]*ModelInfo
}

type ApiInfo struct{
	baseurl		string
	modelname	string
	apikey		string
}

type PromptInfo struct{
	system		string
	user		string
}

type ModelInfo struct{
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
	PromptSystem     string
	PromptUser       string
	ExtraBody        map[string]interface{} `yaml:"extra_body"`
}


//Feed配置
type FeedConfigs struct{
	SouceMap map[string]string
}

func LoadAllConfigs() (){
	sourceMap, _ := LoadFeedConfig("configs/sourceMap.json")

	var Secmindconfigs SecmindConfigs
	Secmindconfigs.feedconfigs.SouceMap = sourceMap



}

func LoadAiConfig(baseDir string) () {
	//先加载model文件
	yamlfile, err := os.ReadFile(filepath.Join(baseDir, "model.yaml"))
    if err != nil {
        return nil, fmt.Errorf("读取 model.yaml 失败: %w", err)
    }

	var modeldata struct {
		Models	[]ModelInfo	`yaml:"models"`
	}
	err = yaml.Unmarshal(yamlfile, &modeldata) 

	//加载api信息
	godotenv.Load(filepath.Join(baseDir, ".env"))
	for _, role := range modeldata.Models{
		role.ModelName = os.Getenv(role.ModelNameEnv)
		role.BaseURL = os.Getenv(role.BaseURLEnv)
		role.APIKey = os.Getenv(role.APIKeyEnv)
	}

	//加载提示词


}

func LoadFeedConfig(source_file_path string) (map[string]string, error) {

	map_file, err := os.Open(source_file_path)
	if err != nil {
		log.Printf("源映射文件加载失败:%s", err)
		return nil, err
	}
	defer map_file.Close()

	map_file_data, err := io.ReadAll(map_file)
	if err != nil {
		log.Printf("从map_file中加载内容失败:%s", err)
		return nil, err
	}

	var sourceMap map[string]string
	err = json.Unmarshal(map_file_data, &sourceMap)

	return sourceMap, err
}

func LoadAllPrompt(promptDir string) (map[string]*PromptInfo, error) {
	dirs, err := os.ReadDir(promptDir)
	if err != nil {
		return nil, err
	}

	promptMap := make(map[string]*PromptInfo)
	
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
	
		sysprompt, _ := os.ReadFile(filepath.Join(promptDir, dir.Name(), "system"))
		usrprompt, _ := os.ReadFile(filepath.Join(promptDir, dir.Name(), "user"))

		promptMap[dir.Name()] = &PromptInfo{
			system: string(sysprompt),
			user: string(usrprompt),
		}
	}
	return promptMap, err
}