package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// RequestsGet 结构体，用于处理 HTTP 请求
type RequestsTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

func NewRequestsTool() *RequestsTool {
	return &RequestsTool{
		Name: "RequestsTool",
		Description: `
		A portal to the internet. Use this when you need to get specific
    content from a website. Input should be a url (i.e. https://www.kubernetes.io/releases).
    The output will be the text response of the GET request.
		`,
		ArgsSchema: `description: "要访问的website，格式是字符串" example: "https://www.kubernetes.io/releases"`,
	}
}

// Run 方法发起 GET 请求并返回处理后的文本
func (r *RequestsTool) Run(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取 URL 失败: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return r.parseHTML(string(body)), nil
}

// parseHTML 方法处理 HTML 内容并返回纯文本
func (r *RequestsTool) parseHTML(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return ""
	}

	// 移除不需要的标签
	doc.Find("header, footer, script, style").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// 获取处理后的纯文本
	return doc.Find("body").Text()
}
