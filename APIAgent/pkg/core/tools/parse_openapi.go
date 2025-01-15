package tools

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/xingyunyang01/APIAgent/pkg/models"
)

// 定义错误类型
var ErrNoServerFound = errors.New("No server found in the swagger yaml.")
var ErrNoPathsFound = errors.New("No paths found in the swagger yaml.")
var ErrNoOperationId = errors.New("No operationId found in operation")
var ErrNoSummaryOrDescription = errors.New("No summary or description found in operation")

type ToolApiSchemaError struct {
	Message string
}

func (e *ToolApiSchemaError) Error() string {
	return e.Message
}

func ParseSwaggerToOpenAPI(swagger *models.Swagger) (*models.OpenAPI, error) {
	// 从 Swagger 获取信息
	info := swagger.Info
	if info == nil {
		info = map[string]string{"title": "Swagger", "description": "Swagger", "version": "1.0.0"}
	}

	servers := swagger.Servers
	if len(servers) == 0 {
		return nil, &ToolApiSchemaError{Message: ErrNoServerFound.Error()}
	}

	// 初始化 OpenAPI 对象
	openAPI := &models.OpenAPI{
		OpenAPI: "3.0.0",
		Info: map[string]string{
			"title":       info["title"],
			"description": info["description"],
			"version":     info["version"],
		},
		Servers:    swagger.Servers,
		Paths:      make(map[string]map[string]interface{}),
		Components: map[string]map[string]interface{}{"schemas": make(map[string]interface{})},
	}

	// 检查 paths 是否存在
	if len(swagger.Paths) == 0 {
		return nil, &ToolApiSchemaError{Message: ErrNoPathsFound.Error()}
	}

	// 转换 paths
	for path, pathItem := range swagger.Paths {
		openAPI.Paths[path] = make(map[string]interface{})
		for method, operation := range pathItem {
			// 将 operation 断言为 map[string]interface{}
			opMap, ok := operation.(map[string]interface{})
			if !ok {
				return nil, &ToolApiSchemaError{Message: fmt.Sprintf("Invalid operation format for %s %s", method, path)}
			}

			operationId, ok := opMap["operationId"].(string)
			if !ok {
				return nil, &ToolApiSchemaError{Message: fmt.Sprintf("%s %s: %s", ErrNoOperationId.Error(), method, path)}
			}

			// 如果没有 summary 或 description, 添加警告
			summary, hasSummary := opMap["summary"].(string)
			description, hasDescription := opMap["description"].(string)

			if !hasSummary && !hasDescription {
				log.Warnf("No summary or description found in operation %s %s.", method, path)
			}

			openAPI.Paths[path][method] = map[string]interface{}{
				"operationId": operationId,
				"summary":     summary,
				"description": description,
				"parameters":  opMap["parameters"],
				"responses":   opMap["responses"],
			}

			// 检查是否存在 requestBody
			if requestBody, ok := opMap["requestBody"]; ok {
				openAPI.Paths[path][method].(map[string]interface{})["requestBody"] = requestBody
			}
		}
	}

	// 转换 definitions
	for name, definition := range swagger.Definitions {
		openAPI.Components["schemas"][name] = definition
	}

	return openAPI, nil
}

// 从 OpenAPI 转为 Tool Bundle
func ParseOpenAPIToToolBundle(openAPI *models.OpenAPI) ([]models.ApiToolBundle, error) {
	serverURL := openAPI.Servers[0]["url"].(string)

	// 列出所有接口
	var interfaces []map[string]interface{}
	for path, pathItem := range openAPI.Paths {
		methods := []string{"get", "post", "put", "delete", "patch", "head", "options", "trace"}
		for _, method := range methods {
			if methodItem, ok := pathItem[method]; ok {
				interfaces = append(interfaces, map[string]interface{}{
					"path":      path,
					"method":    method,
					"operation": methodItem,
				})
			}
		}
	}

	// 获取所有参数并构建工具 bundle
	var bundles []models.ApiToolBundle
	for _, iface := range interfaces {
		parameters := []models.ToolParameter{}
		operation := iface["operation"].(map[string]interface{})

		// 处理参数
		if params, ok := operation["parameters"].([]interface{}); ok {
			for _, param := range params {
				paramMap := param.(map[string]interface{})
				toolParam := models.ToolParameter{
					Name:           paramMap["name"].(string),
					Type:           "string", // 默认类型
					Required:       paramMap["required"].(bool),
					LLMDescription: paramMap["description"].(string),
					Default:        getDefault(paramMap),
				}

				// 类型处理
				if typ := getParameterType(paramMap); typ != "" {
					toolParam.Type = typ
				}

				parameters = append(parameters, toolParam)
			}
		}

		// 处理请求体
		if requestBody, ok := operation["requestBody"].(map[string]interface{}); ok {
			if content, ok := requestBody["content"].(map[string]interface{}); ok {
				for _, contentType := range content {
					if bodySchema, ok := contentType.(map[string]interface{})["schema"].(map[string]interface{}); ok {
						required := bodySchema["required"].([]interface{})
						properties := bodySchema["properties"].(map[string]interface{})
						for name, prop := range properties {
							propMap := prop.(map[string]interface{})
							// 处理引用，如果有
							if ref, ok := propMap["$ref"].(string); ok {
								root := openAPI.Components["schemas"]
								segments := strings.Split(ref, "/")[1:]

								lastSegment := segments[len(segments)-1]
								propMap = root[lastSegment].(map[string]interface{})
							}

							toolParam := models.ToolParameter{
								Name:           name,
								Type:           "string", // 默认类型
								Required:       contains(required, name),
								LLMDescription: propMap["description"].(string),
								Default:        propMap["default"],
							}

							// 如果参数包含 enum，则添加枚举值
							if enum, ok := propMap["enum"].([]interface{}); ok {
								var enumValues []string
								for _, e := range enum {
									enumValues = append(enumValues, e.(string))
								}
								toolParam.Enum = enumValues
							}

							// 类型处理
							if typ := getParameterType(propMap); typ != "" {
								toolParam.Type = typ
							}

							parameters = append(parameters, toolParam)
						}
					}
				}
			}
		}

		// 检查参数是否重复
		paramCount := make(map[string]int)
		for _, param := range parameters {
			paramCount[param.Name]++
		}
		for name, count := range paramCount {
			if count > 1 {
				log.Warnf("Parameter %s is duplicated.", name)
			}
		}

		// 设置 operationId
		if _, ok := operation["operationId"]; !ok {
			// 如果没有 operationId，使用 path 和 method 生成
			path := iface["path"].(string)
			if strings.HasPrefix(path, "/") {
				path = path[1:]
			}
			// 移除特殊字符以确保 operationId 合法
			re := regexp.MustCompile("[^a-zA-Z0-9_-]")
			path = re.ReplaceAllString(path, "")
			if path == "" {
				path = uuid.New().String()
			}
			operation["operationId"] = fmt.Sprintf("%s_%s", path, iface["method"].(string))
		}

		// 构建 ApiToolBundle
		bundles = append(bundles, models.ApiToolBundle{
			ServerURL:   serverURL + iface["path"].(string),
			Method:      iface["method"].(string),
			Summary:     getStringOrDefault(operation["description"], ""),
			OperationID: operation["operationId"].(string),
			Parameters:  parameters,
			OpenAPI:     operation,
		})
	}

	return bundles, nil
}

func getStringOrDefault(value interface{}, defaultValue string) string {
	if value == nil {
		return defaultValue
	}
	strValue, ok := value.(string)
	if !ok {
		return defaultValue
	}
	return strValue
}

// 获取默认值
func getDefault(param map[string]interface{}) interface{} {
	if schema, ok := param["schema"].(map[string]interface{}); ok {
		return schema["default"]
	}
	return nil
}

// 获取参数类型
func getParameterType(param map[string]interface{}) string {
	// 根据实际情况来返回正确的类型
	if schema, ok := param["schema"].(map[string]interface{}); ok {
		if typ, ok := schema["type"].(string); ok {
			return typ
		}
	}
	return ""
}

// 判断元素是否在数组中
func contains(slice []interface{}, elem string) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}
	return false
}
