package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/xingyunyang01/agent/pkg/ai"
	prompttpl "github.com/xingyunyang01/agent/pkg/promptTpl"
	"github.com/xingyunyang01/agent/pkg/tools"
)

func main() {
	query := "1+2+3+4-5-6=? Just give me a number result"

	addtool := tools.AddToolName + ":" + tools.AddToolDescription + "\nparam: \n" + tools.AddToolParam
	subtool := tools.SubToolName + ":" + tools.SubToolDescription + "\nparam: \n" + tools.SubToolParam
	toolsL := make([]string, 0)
	toolsL = append(toolsL, addtool, subtool)

	tool_names := make([]string, 0)
	tool_names = append(tool_names, tools.AddToolName, tools.SubToolName)

	prompt := fmt.Sprintf(prompttpl.Template, toolsL, tool_names, query)
	fmt.Println("prompt: ", prompt)
	//注入用户prompt
	ai.MessageStore.AddForUser(prompt)
	i := 1
	for {
		first_response := ai.NormalChat(ai.MessageStore.ToMessage())
		fmt.Printf("========第%d轮回答========\n", i)
		fmt.Println(first_response)
		regexPattern := regexp.MustCompile(`Final Answer:\s*(.*)`)
		finalAnswer := regexPattern.FindStringSubmatch(first_response.Content)
		if len(finalAnswer) > 1 {
			fmt.Println("========最终 GPT 回复========")
			fmt.Println(first_response.Content)
			break
		}

		ai.MessageStore.AddForAssistant(first_response)

		regexAction := regexp.MustCompile(`Action:\s*(.*?)[.\n]`)
		regexActionInput := regexp.MustCompile(`Action Input:\s*(.*?)[.\n]`)

		action := regexAction.FindStringSubmatch(first_response.Content)
		actionInput := regexActionInput.FindStringSubmatch(first_response.Content)

		if len(action) > 1 && len(actionInput) > 1 {
			i++
			result := 0
			//需要调用工具
			if action[1] == "AddTool" {
				fmt.Println("calls AddTool")
				result = tools.AddTool(actionInput[1])
			} else if action[1] == "SubTool" {
				fmt.Println("calls SubTool")
				result = tools.SubTool(actionInput[1])
			}
			fmt.Println("========函数返回结果========")
			fmt.Println(result)

			Observation := "Observation: " + strconv.Itoa(result)
			prompt = first_response.Content + Observation
			fmt.Printf("========第%d轮的prompt========\n", i)
			fmt.Println(prompt)
			ai.MessageStore.AddForUser(prompt)
		}
	}
}
