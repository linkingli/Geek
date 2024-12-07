package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xingyunyang01/ginTools/pkg/config"
	"github.com/xingyunyang01/ginTools/pkg/controllers"
	"github.com/xingyunyang01/ginTools/pkg/services"
)

func main() {
	k8sconfig := config.NewK8sConfig().InitRestConfig()
	restMapper := k8sconfig.InitRestMapper()
	dynamicClient := k8sconfig.InitDynamicClient()
	informer := k8sconfig.InitInformer()

	clientSet := k8sconfig.InitClientSet()

	resourceCtl := controllers.NewResourceCtl(services.NewResourceService(&restMapper, dynamicClient, informer))
	podLogCtl := controllers.NewPodLogEventCtl(services.NewPodLogEventService(clientSet))

	r := gin.New()

	r.GET("/:resource", resourceCtl.List())
	r.DELETE("/:resource", resourceCtl.Delete())
	r.POST("/:resource", resourceCtl.Create())
	r.GET("/pods/logs", podLogCtl.GetLog())
	r.GET("/pods/events", podLogCtl.GetEvent())

	r.Run(":8080")
}
