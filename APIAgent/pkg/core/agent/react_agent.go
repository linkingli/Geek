package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	outputparser "github.com/xingyunyang01/APIAgent/pkg/core/agent/output_parser"
	promptTemplate "github.com/xingyunyang01/APIAgent/pkg/core/agent/template"
	"github.com/xingyunyang01/APIAgent/pkg/core/ai"
	"github.com/xingyunyang01/APIAgent/pkg/core/tools"
	"github.com/xingyunyang01/APIAgent/pkg/models"
)

func organizeToolsPrompt(tools []models.ApiToolBundle) (string, []string) {
	messageTools := make([]models.PromptMessageTool, 0)
	toolNames := []string{}
	for _, tool := range tools {
		toolNames = append(toolNames, tool.OperationID)

		properties := make([]models.Properties, 0)
		for _, parameter := range tool.Parameters {
			properties = append(properties, models.Properties{
				Type:          parameter.Type,
				Name:          parameter.Name,
				Descripeption: parameter.LLMDescription,
				Enum:          parameter.Enum,
			})
		}

		required := make([]string, 0)
		for _, parameter := range tool.Parameters {
			if parameter.Required {
				required = append(required, parameter.Name)
			}
		}

		parameters := models.Parameters{
			Type:       "object",
			Properties: properties,
			Required:   required,
		}

		messageTool := models.PromptMessageTool{
			Name:        tool.OperationID,
			Description: tool.Summary,
			Parameters:  parameters,
		}

		messageTools = append(messageTools, messageTool)
	}

	jsonMessageTools, err := json.Marshal(messageTools)
	if err != nil {
		panic(err)
	}

	//fmt.Println("json: ", string(jsonMessageTools))

	return string(jsonMessageTools), toolNames
}

func organizeHistoricMessage() string {
	return ""
}

func organizeReActTemplate(instruction string, tools []models.ApiToolBundle, query string) string {
	messageTools, messageToolsNames := organizeToolsPrompt(tools)

	historicMessage := organizeHistoricMessage()

	// 填充数据
	data := models.TemplateData{
		Instruction:      instruction,
		Tools:            messageTools,
		ToolNames:        messageToolsNames,
		HistoricMessages: historicMessage,
		Query:            query,
	}

	// Load and render the templatet
	tmpl, err := template.New("prompt").Parse(promptTemplate.EN_Template)
	if err != nil {
		log.Fatalln("Failed to parse template:", err)
	}

	// 创建一个缓冲区来接收模板的渲染输出
	var result bytes.Buffer

	// 执行模板并填充数据
	err = tmpl.Execute(&result, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return result.String()
}

func Run(sc *models.Config, toolBundles []models.ApiToolBundle, query string) (string, error) {
	prompt := organizeReActTemplate(sc.Instruction, toolBundles, query)
	ai.MessageStore.AddForUser(prompt)

	var action string
	var actionInput map[string]interface{}
	var err error

	iteration_steps := 1
	for {
		first_response := ai.NormalChat(ai.MessageStore.ToMessage())
		fmt.Printf("========第%d轮回答========\n", iteration_steps)
		fmt.Println(first_response.Content)
		ai.MessageStore.AddForAssistant(first_response.Content)

		action, actionInput, err = outputparser.HandleReActOutput(first_response.Content)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}

		fmt.Println("Action:", action)
		fmt.Println("Action_Input:", actionInput)

		if action == "Final Answer" {
			break
		}

		if iteration_steps >= sc.MaxIterationSteps {
			actionInput["finalAnswer"] = "已超出允许的最大迭代次数"
			break
		}

		iteration_steps++
		Observation := "Observation: %s"
		for _, toolBundle := range toolBundles {
			if action == toolBundle.OperationID {
				respBody, statusCode, err := tools.ToolInvoke(sc.APIs.APIProvider.APIKey, toolBundle.Method, toolBundle.ServerURL, toolBundle, actionInput)
				if err != nil {
					fmt.Println("Error:", err)
					return "", err
				}
				fmt.Println("StatusCode:", statusCode)
				fmt.Println("Response:", string(respBody))
				Observation = fmt.Sprintf(Observation, string(respBody))
				break
			}
		}
		ai.MessageStore.AddForUser(Observation)
	}
	return actionInput["finalAnswer"].(string), nil
}
