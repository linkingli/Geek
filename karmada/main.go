package main

import (
	"context"
	"fmt"

	"github.com/xingyunyang01/karmada/pkg/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func test() {
	config := config.NewK8sConfig().InitRestConfig()

	ri := config.InitDynamicClient()

	clusterGvr := schema.GroupVersionResource{
		Group:    "cluster.karmada.io",
		Version:  "v1alpha1",
		Resource: "clusters",
	}
	list, err := ri.Resource(clusterGvr).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, item := range list.Items {
		fmt.Println(item.GetName())
	}
}

func main() {
	config := config.NewK8sConfig().InitRestConfig()

	client := config.InitClientSet()

	clusters, err := client.ClusterV1alpha1().Clusters().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, item := range clusters.Items {
		fmt.Println(item.GetName())
	}
}
