package main

import (
	"encoding/json"
	"net/http"

	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

func main() {
	wrapper.SetCtx(
		// 插件名称
		"my-plugin",
		// 为解析插件配置，设置自定义函数
		wrapper.ParseConfigBy(parseConfig),
		// 为处理返回体，设置自定义函数
		wrapper.ProcessResponseBodyBy(onHttpResponseBody),
	)
}

// completion
type Completion struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Deepl struct {
	Text        []string `json:"text"`
	Target_lang string   `json:"target_lang"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Choices           []Choice        `json:"choices"`
	Object            string          `json:"object"`
	Usage             CompletionUsage `json:"usage"`
	Created           string          `json:"created"`
	SystemFingerprint string          `json:"system_fingerprint"`
	Model             string          `json:"model"`
	ID                string          `json:"id"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// 自定义插件配置
type PluginConfig struct {
	url       string
	model     string
	apiKey    string
	LLMClient wrapper.HttpClient
}

// 在控制台插件配置中填写的YAML配置会自动转换为JSON，此处直接从JSON这个参数里解析配置即可
func parseConfig(json gjson.Result, config *PluginConfig, log wrapper.Log) error {
	log.Info("[parseConfig] start")
	// 解析出配置，更新到config中
	config.url = json.Get("url").String()
	config.model = json.Get("model").String()
	config.apiKey = json.Get("apiKey").String()

	config.LLMClient = wrapper.NewClusterClient(wrapper.FQDNCluster{
		FQDN: json.Get("serviceFQDN").String(),
		Port: json.Get("servicePort").Int(),
		Host: json.Get("serviceHost").String(),
	})
	return nil
}

// 从response接收到firstreq的大模型返回
func onHttpResponseBody(ctx wrapper.HttpContext, config PluginConfig, body []byte, log wrapper.Log) types.Action {
	var responseCompletion CompletionResponse
	_ = json.Unmarshal(body, &responseCompletion)
	log.Infof("content: %s", responseCompletion.Choices[0].Message.Content)

	completion := Completion{
		Model: config.model,
		Messages: []Message{{Role: "system", Content: `请参考我的 JSON 定义输出 JSON 对象，示例：{"ouput": "xxxx"}`},
			{Role: "user", Content: responseCompletion.Choices[0].Message.Content}},
		Stream: false,
	}
	headers := [][2]string{{"Content-Type", "application/json"}, {"Authorization", "Bearer " + config.apiKey}}
	reqEmbeddingSerialized, _ := json.Marshal(completion)
	err := config.LLMClient.Post(
		config.url,
		headers,
		reqEmbeddingSerialized,
		func(statusCode int, responseHeaders http.Header, responseBody []byte) {
			log.Infof("statusCode: %d", statusCode)
			log.Infof("responseBody: %s", string(responseBody))
			//得到gpt的返回结果
			var responseCompletion CompletionResponse
			_ = json.Unmarshal(responseBody, &responseCompletion)
			log.Infof("content: %s", responseCompletion.Choices[0].Message.Content)

			if responseCompletion.Choices[0].Message.Content != "" {
				//如果结果不是空，则替换原本的response body
				newbody, err := json.Marshal(responseCompletion.Choices[0].Message.Content)
				if err != nil {
					proxywasm.ResumeHttpResponse()
					return
				}
				proxywasm.ReplaceHttpResponseBody(newbody)
				proxywasm.ResumeHttpResponse()
			}
			log.Infof("resume")
			proxywasm.ResumeHttpResponse()
		}, 50000)
	if err != nil {
		log.Errorf("[onHttpResponseBody] completion err: %s", err.Error())
		proxywasm.ResumeHttpResponse()
	}
	return types.ActionPause
}
