package sys

import (
	"log"

	"github.com/xingyunyang01/APIAgent/pkg/models"
	"gopkg.in/yaml.v3"
)

func newSysConfig() *models.Config {
	apis := models.APIConfig{}
	return &models.Config{APIs: apis}
}

func InitConfig() *models.Config {
	config := newSysConfig()

	if b := LoadConfigFile(); b != nil { //读取 system.yaml文件内容
		err := yaml.Unmarshal(b, config) //将byte反序列化成结构体
		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}
