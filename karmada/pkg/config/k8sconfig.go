package config

import (
	"path/filepath"

	"github.com/pkg/errors"

	karmadaversiond "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sConfig struct {
	*rest.Config
	*karmadaversiond.Clientset
	*dynamic.DynamicClient
	e error
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

// 初始化k8s配置
func (k *K8sConfig) InitRestConfig() *K8sConfig {
	kuebconfig := filepath.Join(homedir.HomeDir(), ".kube", "karmada-config")
	config, _ := clientcmd.BuildConfigFromFlags("", kuebconfig)
	k.Config = config
	return k
}

func (k *K8sConfig) Error() error {
	return k.e
}

// 初始化clientSet客户端
func (k *K8sConfig) InitClientSet() *karmadaversiond.Clientset {
	if k.Config == nil {
		k.e = errors.Wrap(errors.New("k8s config is nil"), "init k8s client failed")
		return nil
	}

	clientSet, err := karmadaversiond.NewForConfig(k.Config)
	if err != nil {
		k.e = errors.Wrap(err, "init karmada clientSet failed")
		return nil
	}
	return clientSet
}

// 初始化动态客户端
func (k *K8sConfig) InitDynamicClient() *dynamic.DynamicClient {
	if k.Config == nil {
		k.e = errors.Wrap(errors.New("k8s config is nil"), "init k8s client failed")
		return nil
	}

	dynamicClient, err := dynamic.NewForConfig(k.Config)
	if err != nil {
		k.e = errors.Wrap(err, "init k8s dynamicClient failed")
		return nil
	}
	return dynamicClient
}
