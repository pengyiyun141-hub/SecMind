package analyzer
/*
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


func NewClient(role string)(*Client, error){
	//现根据role加载对应的配置清单
	ai_configmanifest_filePath := fmt.Sprintf("configs/%s_aiconfig_manifest.yaml", role) 
	ai_configmanifest_file, err:= os.ReadFile(ai_configmanifest_filePath)

	var test struct{
		configManifest ConfigPath `yaml:"configpath"`
	}
	err = yaml.Unmarshal(ai_configmanifest_file, &test) 

	//将api信息加载到环境变量
	err = godotenv.Load(test.configManifest.EnvFile)
	if err != nil {
		log.Println("加载 .env 失败: ", err)
	}

	//加载yaml配置文件
	var all_model_param *[]ModelSpec
	all_model_param, err = loadModelYamlFile(test.configManifest.ModelCofig)
	if err != nil {
		log.Fatal("加载 .env 失败: ", err)
	}
	fmt.Print(all_model_param)
	//return err
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
}*/