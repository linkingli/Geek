package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// getHTTP 执行一个GET HTTP请求到指定的URL，并打印响应体。
func GetHTTP(url string) (string, error) {
	// 创建HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func PostHTTP(url string, body []byte) (string, error) {
	// 创建HTTP POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// DeleteHTTP 执行一个DELETE HTTP请求到指定的URL，并打印响应体。
func DeleteHTTP(url string) (string, error) {
	// 创建HTTP DELETE请求
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
