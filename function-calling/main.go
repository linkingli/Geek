package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xingyunyang01/Geek/function-calling/ai"
	"github.com/xingyunyang01/Geek/function-calling/tools"

	"github.com/sashabaranov/go-openai"
)

func main() {
	toolsList := make([]openai.Tool, 0)
	toolsList = append(toolsList, tools.AddToolDefine, tools.SubToolDefine)

	prompt := "1+2=? Just give me a number result"
	ai.MessageStore.AddFor(ai.RoleUser, prompt, nil)

	response := ai.ToolsChat(ai.MessageStore.ToMessage(), toolsList)
	toolCalls := response.ToolCalls

	fmt.Println("大模型的回复是: ", response.Content)
	fmt.Println("大模型选择的工具是: ", toolCalls)

	if toolCalls != nil {
		var result int
		var args tools.InputArgs
		err := json.Unmarshal([]byte(toolCalls[0].Function.Arguments), &args)
		if err != nil {
			log.Fatalln("json unmarshal err: ", err.Error())
		}

		if toolCalls[0].Function.Name == tools.AddToolDefine.Function.Name {
			result = tools.AddTool(args.Numbers)
		} else if toolCalls[0].Function.Name == tools.SubToolDefine.Function.Name {
			result = tools.SubTool(args.Numbers)
		}

		fmt.Println("函数计算结果: ", result)
		ai.MessageStore.AddFor(ai.RoleAssistant, response.Content, toolCalls)
		ai.MessageStore.AddForTool(string(result), toolCalls[0].Function.Name, toolCalls[0].ID)

		response := ai.ToolsChat(ai.MessageStore.ToMessage(), toolsList)

		fmt.Println("大模型的最终回复: ", response.Content)
	}
}
