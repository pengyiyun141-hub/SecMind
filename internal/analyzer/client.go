package analyzer

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Client struct {
	all_model_info []ModelSpec
}

func (c *Client) Execute(templatestName string, userInput string)(string, error){
	
}


func NewClient(configPath string)(*Client, error){
	envPath := configPath + "/.env"
	err := godotenv.Load(envPath)
	if err != nil {
		log.Println("加载 .env 失败: ", err)
	}

	//加载yaml配置文件
	var all_model_param *[]ModelSpec
	modelPath := configPath + "/model.yaml"
	all_model_param, err = loadModelYamlFile(modelPath)
	if err != nil {
		log.Fatal("加载 .env 失败: ", err)
	}
	fmt.Print(all_model_param)
	return err
}

func loadModelYamlFile(configYamlFilePath string)(*[]ModelSpec, error){
	yamlfile, err := os.ReadFile(configYamlFilePath)
	if err != nil {
		return nil, fmt.Errorf("加载%s文件失败: %w", configYamlFilePath, err)
	}

	var wrapper struct {
		Models []ModelSpec `yaml:"models"`
	}

	err = yaml.Unmarshal(yamlfile, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("解析yaml值失败: %w", err)
	}

	return &wrapper.Models, err
}