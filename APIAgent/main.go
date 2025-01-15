package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xingyunyang01/APIAgent/pkg/controllers"
	"github.com/xingyunyang01/APIAgent/pkg/core/tools"
	"github.com/xingyunyang01/APIAgent/pkg/models"
	"github.com/xingyunyang01/APIAgent/pkg/services"
	"github.com/xingyunyang01/APIAgent/pkg/sys"
	"gopkg.in/yaml.v3"
)

func main() {
	sc := sys.InitConfig()
	//fmt.Printf("%+v\n", sc)
	var api models.OpenAPI
	err := yaml.Unmarshal([]byte(sc.APIs.API), &api)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("%+v\n", api)

	tools, err := tools.ParseOpenAPIToToolBundle(&api)

	chatCompletionService := services.NewChatCompletionService(sc, tools)

	chatCompletionCtl := controllers.NewChatCompletionCtl(chatCompletionService)

	r := gin.New()

	r.POST("/v1/chat-messages", chatCompletionCtl.ChatCompletion())

	r.Run(":8080")

	//agent.Run(sc, tools, `帮我把"何以解忧，唯有暴富"翻译成英文`)
	//agent.Run(sc, tools, `济南奥体中心附近游泳馆有哪些`)
}
