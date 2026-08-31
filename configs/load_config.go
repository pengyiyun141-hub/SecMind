package configs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// LoadAllConfigs()是基础模块，其执行失败则整个程序没有往后执行的必要。
func LoadAllConfigs() (*SecmindConfigs, error) {
	SecCfgs := &SecmindConfigs{
		Feedconfigs: &FeedConfigs{},
		Aiconfigs:   nil,
	}

	var err error

	FeedConfigsSourceInfoFilesPath := filepath.Join("configs", "sourceinfo", "*.json")
	FeedConfigsFiles, err := filepath.Glob(FeedConfigsSourceInfoFilesPath)
	if err != nil {
		return nil, fmt.Errorf("LoadFeedConfig(),FeedConfigsFiles文件获取失败：%w\n", err)
	}

	sourceinfomap := make(map[string]*SourceInfo)
	for _, FeedConfigsFile := range FeedConfigsFiles {
		sourceinfo, err := LoadFeedConfig(FeedConfigsFile)
		if err != nil {
			return nil, fmt.Errorf("文件%s:LoadFeedConfig()执行失败：%w\n", FeedConfigsFile, err)
		}
		sourceinfomap[sourceinfo.SourceName] = sourceinfo
	}
	SecCfgs.Feedconfigs.SourceInfoMap = sourceinfomap

	SecCfgs.Aiconfigs, err = LoadAiConfig(filepath.Join("configs"))
	if err != nil {
		return nil, fmt.Errorf("LoadAiConfig()执行失败：%w\n", err)
	}

	SecCfgs.Poolconfigs, err = LoadPoolConfig(filepath.Join("configs", "pool.json"))
	if err != nil {
		return nil, fmt.Errorf("LoadPoolConfig()执行失败：%w\n", err)
	}

	SecCfgs.Logconfigs, err = LoadSecMindLog(filepath.Join("configs", "secmindlog.yaml"))
	if err != nil {
		return nil, fmt.Errorf("LoadSecMindLog()执行失败：%w\n", err)
	}
	
	return SecCfgs, err
}

func LoadAiConfig(baseDir string) (*AiConfigs, error) {
	airole := &AiConfigs{
		Modelinfo:  make(map[string]*ModelInfo),
		Promptinfo: make(map[string]*PromptInfo),
		Apiinfo:    make(map[string]*ApiInfo),
	}
	//先加载model文件
	yamlFileData, err := os.ReadFile(filepath.Join(baseDir, "model.yaml"))
	if err != nil {
		return nil, fmt.Errorf("读取 model.yaml 失败: %w", err)
	}

	var modeldata struct {
		ModelCfgs []ModelInfo `yaml:"models"`
	}
	err = yaml.Unmarshal(yamlFileData, &modeldata)
	//fmt.Printf("yaml:%s", string(yamlfile))

	//加载api信息
	godotenv.Load(filepath.Join(baseDir, ".env"))
	for i := range modeldata.ModelCfgs {
		role := &modeldata.ModelCfgs[i]
		role.ModelName = os.Getenv(role.ModelNameEnv)
		role.BaseURL = os.Getenv(role.BaseURLEnv)
		role.APIKey = os.Getenv(role.APIKeyEnv)
		airole.Modelinfo[role.Name] = role
	}

	//加载提示词
	airole.Promptinfo, err = LoadAllPrompt("configs/prompts/")
	return airole, err
}

func LoadFeedConfig(sourceInfoFilePath string) (*SourceInfo, error) {

	map_file, err := os.Open(sourceInfoFilePath)
	if err != nil {
		return nil, fmt.Errorf("源映射文件加载失败:%w", err)
	}
	defer map_file.Close()

	map_file_data, err := io.ReadAll(map_file)
	if err != nil {
		return nil, fmt.Errorf("从map_file中加载内容失败:%w", err)
	}

	sourceinfo := &SourceInfo{}
	err = json.Unmarshal(map_file_data, sourceinfo)

	return sourceinfo, err
}

func LoadPoolConfig(poolConfigPath string) (*PoolConfigs, error) {
	configdata, err := os.ReadFile(poolConfigPath)
	if err != nil {
		return nil, fmt.Errorf("[err]poolConfig文件加载失败:%w", err)
	}

	poolconfig := &PoolConfigs{}
	err = json.Unmarshal(configdata, &poolconfig)
	if err != nil {
		return nil, fmt.Errorf("[err]poolConfig解析失败:%w", err)
	}

	return poolconfig, nil
}

func LoadAllPrompt(promptDir string) (map[string]*PromptInfo, error) {
	dirs, err := os.ReadDir(promptDir)
	if err != nil {
		return nil, fmt.Errorf("加载提示词文件夹失败：%w", err)
	}

	promptMap := make(map[string]*PromptInfo)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		sysprompt, _ := os.ReadFile(filepath.Join(promptDir, dir.Name(), "system.txt"))
		usrprompt, _ := os.ReadFile(filepath.Join(promptDir, dir.Name(), "user.txt"))

		promptMap[dir.Name()] = &PromptInfo{
			System: string(sysprompt),
			User:   string(usrprompt),
		}
	}
	return promptMap, err
}

func LoadSecMindLog(OptionsPath string) (*LogConfigs, error) {
	logconfigdata := &LogConfigs{}

	logConfigFile, err := os.ReadFile(OptionsPath)
	if err != nil {
		return nil, fmt.Errorf("加载secmindlog.yaml文件路径失败：%w", err)
	}

	err = yaml.Unmarshal(logConfigFile, &logconfigdata)
	if err != nil {
		return nil, fmt.Errorf("加载secmindlog.yaml文件失败：%w", err)
	}

	return logconfigdata, nil
}
