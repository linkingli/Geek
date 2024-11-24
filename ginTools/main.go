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

	ctl := controllers.NewResourceCtl(services.NewResourceService(&restMapper, dynamicClient, informer))

	r := gin.New()

	r.GET("/:resource", ctl.List())
	r.DELETE("/:resource", ctl.Delete())
	r.POST("/:resource", ctl.Create())

	r.Run(":8080")
}
