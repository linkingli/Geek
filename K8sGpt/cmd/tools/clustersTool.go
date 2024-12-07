package tools

import (
	"github.com/xingyunyang01/K8sGpt/cmd/utils"
)

// ListTool 表示一个工具，用于列出 k8s 资源命令。
type ClusterTool struct {
	Name        string
	Description string
}

// NewListTool 创建一个新的 ListTool 实例。
func NewClusterTool() *ClusterTool {
	return &ClusterTool{
		Name:        "ClusterTool",
		Description: "用于列出集群列表",
	}
}

// Run 执行命令并返回输出。
func (l *ClusterTool) Run() (string, error) {

	url := "http://localhost:8081/clusters"

	s, err := utils.GetHTTP(url)

	return s, err
}
