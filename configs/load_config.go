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
	Aiconfigs		*AiConfigs
	Feedconfigs     *FeedConfigs 
}

//AI配置
type AiConfigs struct{
	Apiinfo		map[string]*ApiInfo
	Promptinfo	map[string]*PromptInfo
	Modelinfo	map[string]*ModelInfo
}

type ApiInfo struct{
	Baseurl		string
	Modelname	string
	Apikey		string
}

type PromptInfo struct{
	System		string
	User		string
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

//LoadAllConfigs()是基础模块，其执行失败则整个程序没有往后执行的必要。
func LoadAllConfigs() (*SecmindConfigs, error){
	fmt.Println("执行 LoadAllConfigs()")
	var Secmindconfigs struct {
		Aiconfigs   *AiConfigs
    	Feedconfigs *FeedConfigs
	}
	var err error
	
	fmt.Println("执行 LoadFeedConfig()")
	Secmindconfigs.Feedconfigs.SouceMap, err = LoadFeedConfig("configs/sourceMap.json")
	if err != nil {
		log.Fatal("sourceMap.json文件加载失败")
	}
	fmt.Println("执行 LoadAllConfigs()失败")

	Secmindconfigs.Aiconfigs, err = LoadAiConfig("configs/")

	fmt.Printf("%s/n", Secmindconfigs.Feedconfigs.SouceMap)

	return &Secmindconfigs, err
}

func LoadAiConfig(baseDir string) (*AiConfigs, error) {
	airole := &AiConfigs {
		Modelinfo:  make(map[string]*ModelInfo),
        Promptinfo: make(map[string]*PromptInfo),
        Apiinfo:    make(map[string]*ApiInfo),
	}
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
	for i := range modeldata.Models{
		role := &modeldata.Models[i]
		role.ModelName = os.Getenv(role.ModelNameEnv)
		role.BaseURL = os.Getenv(role.BaseURLEnv)
		role.APIKey = os.Getenv(role.APIKeyEnv)
		airole.Modelinfo[role.Name] = role
	}

	//加载提示词
	airole.Promptinfo, err = LoadAllPrompt("configs/prompts/")
	return airole, err
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
			System: string(sysprompt),
			User: string(usrprompt),
		}
	}
	return promptMap, err
}