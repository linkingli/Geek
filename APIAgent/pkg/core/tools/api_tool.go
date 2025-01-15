package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/xingyunyang01/APIAgent/pkg/models"
)

func assembingRequest(apiKey models.APIKey, url string) (map[string]string, string) {
	if apiKey.In == "header" {
		headers := make(map[string]string, 2)
		headers["Content-Type"] = "application/json"
		headers["Authorization"] = apiKey.Name + " " + apiKey.Value

		return headers, url
	} else if apiKey.In == "query" {
		url += "?" + apiKey.Name + "=" + apiKey.Value
		return nil, url
	} else {
		return nil, url
	}
}

func doHttpRequest(method, urlStr string, headers map[string]string, toolBundle models.ApiToolBundle, actionInput map[string]interface{}) ([]byte, int, error) {
	var reqBody []byte
	var err error

	// 解析URL模板以查找路径参数
	urlParts := strings.Split(urlStr, "/")
	for i, part := range urlParts {
		if strings.Contains(part, "{") && strings.Contains(part, "}") {
			for _, param := range toolBundle.Parameters {
				paramNameInPath := part[1 : len(part)-1]
				if paramNameInPath == param.Name {
					if value, ok := actionInput[param.Name]; ok {
						// 删除已经使用过的
						delete(actionInput, param.Name)
						// 替换模板中的占位符
						urlParts[i] = url.QueryEscape(value.(string))
					}
				}
			}
		}
	}

	urlStr = strings.Join(urlParts, "/")

	if toolBundle.OpenAPI["requestBody"] != nil {
		reqBody, err = json.Marshal(actionInput)
		if err != nil {
			return nil, 400, err
		}
	} else {
		reqBody = nil
		for _, param := range toolBundle.Parameters {
			urlStr += "&" + param.Name + "=" + actionInput[param.Name].(string)
		}
	}

	fmt.Println("method: ", method)
	fmt.Println("urlStr: ", urlStr)
	fmt.Println("headers: ", headers)
	fmt.Println("reqBody: ", string(reqBody))

	return call(method, urlStr, headers, reqBody)
}

func call(method, url string, headers map[string]string, reqBody []byte) ([]byte, int, error) {
	method = strings.ToUpper(method)
	// 创建请求体
	var body *bytes.Reader
	if reqBody != nil {
		body = bytes.NewReader(reqBody)
	} else {
		body = bytes.NewReader([]byte{})
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 创建 HTTP 客户端，并设置超时时间
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("读取响应体失败: %v", err)
	}

	return respBody, resp.StatusCode, nil
}

func ToolInvoke(apiKey models.APIKey, method, url string, toolBundle models.ApiToolBundle, actionInput map[string]interface{}) ([]byte, int, error) {
	headers, url := assembingRequest(apiKey, url)

	respBody, statusCode, err := doHttpRequest(method, url, headers, toolBundle, actionInput)

	return respBody, statusCode, err
}
