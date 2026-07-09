package analyzer

import (
	"fmt"
	"io"
	"os"
)

//由于在analyzer模块中的每个文件结构都极其相似，相似的部分为开头的一部分加载配置外加发送请求。因此重构计划预计将这段内容打包为一个client文件，该文件专注于加载配置和发送请求。
func AnalyzeArticleByAi(model_param *ModelSpec, article_Path string) (string, error){
	articleContextFile, err := os.Open(article_Path)
	if err != nil {
		return "", fmt.Errorf("打开文章失败")
	}
	
	articleContext, err:= io.ReadAll(articleContextFile)
	if err != nil || len(articleContext) == 0 {
		return "", fmt.Errorf("文章数据读入内存失败")
	}
	//fmt.Printf("\n打开文件获得的byte内容为：%s\n", articleContext)

	var promptSys string
	var promptText string
	promptTextdata_sys, err := os.ReadFile(model_param.SystemPrompt)
	if err != nil {
		return "", fmt.Errorf("打开%s文件失败", model_param.SystemPrompt)
	}
	promptTextdata_user, err := os.ReadFile(model_param.UserPrompt)
	if err != nil {
		return "", fmt.Errorf("打开%s文件失败", model_param.UserPrompt)
	}

	promptSys = string(promptTextdata_sys)
	promptText = string(promptTextdata_user)

	promptText += fmt.Sprintf("%s", string(articleContext))
	
	var promptMessage []Message
	promptMessage = []Message{
		{Role: "system", Content: promptSys},
		{Role: "user", Content: promptText},
	}


	//fmt.Printf("\n准备发给AI的文章内容为：\n%s", promptText)
	text, err:= CallAiApi(model_param, promptMessage)
	

	return string(text), err
}
