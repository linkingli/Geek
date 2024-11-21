package tools

import (
	"strconv"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func AddToolTpl() openai.Tool {
	fundefine := openai.FunctionDefinition{
		Name:        AddToolName,
		Description: AddToolDescription,
		Parameters:  AddToolParam,
	}

	tool := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &fundefine,
	}

	return tool
}

func SubToolTpl() openai.Tool {
	fundefine := openai.FunctionDefinition{
		Name:        SubToolName,
		Description: SubToolDescription,
		Parameters:  SubToolParam,
	}

	tool := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &fundefine,
	}

	return tool
}

func ToolsTpl(name string, description string, param string, tools []openai.Tool) []openai.Tool {
	fundefine := openai.FunctionDefinition{
		Name:        name,
		Description: description,
		Parameters:  param,
	}

	tool := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &fundefine,
	}

	tools = append(tools, tool)

	return tools
}

func AddTool(numbers string) int {
	num := strings.Split(numbers, ",")
	inum0, _ := strconv.Atoi(num[0])
	inum1, _ := strconv.Atoi(num[1])
	return inum0 + inum1
}

func SubTool(numbers string) int {
	num := strings.Split(numbers, ",")
	inum0, _ := strconv.Atoi(num[0])
	inum1, _ := strconv.Atoi(num[1])
	return inum0 - inum1
}
