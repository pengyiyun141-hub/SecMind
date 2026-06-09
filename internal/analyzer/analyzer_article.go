package analyzer

import (
	"fmt"
	"io"
	"os"
)

func AnalyzeArticleByAi(model_param *ModelSpec) (string, error){

	articleContextFile, err := os.Open("internal/data/articles/Subtle Injection for Ground-truth Inference of LLM Training Data")
	if err != nil {
		return "", fmt.Errorf("打开文章失败")
	}
	
	articleContext, err:= io.ReadAll(articleContextFile)
	if err != nil || len(articleContext) == 0 {
		return "", fmt.Errorf("文章数据读入内存失败")
	}

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



	text, err:= CallAiApi(model_param, promptMessage)
	
	return string(text), err
}
