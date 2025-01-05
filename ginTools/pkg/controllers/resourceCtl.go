package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xingyunyang01/ginTools/pkg/services"
)

type ResourceCtl struct {
	resourceService *services.ResourceService
}

func NewResourceCtl(service *services.ResourceService) *ResourceCtl {
	return &ResourceCtl{resourceService: service}
}

func (r *ResourceCtl) List() func(c *gin.Context) {
	return func(c *gin.Context) {
		var resource = c.Param("resource")
		ns := c.DefaultQuery("ns", "default")
		resourceList, _ := r.resourceService.ListResource(resource, ns)
		c.JSON(200, gin.H{"data": resourceList})
	}
}

func (r *ResourceCtl) Delete() func(c *gin.Context) {
	return func(c *gin.Context) {
		var resource = c.Param("resource")
		ns := c.DefaultQuery("ns", "default")
		name := c.Query("name")
		err := r.resourceService.DeleteResource(resource, ns, name)
		if err != nil {
			c.JSON(500, gin.H{"error": "删除失败：" + err.Error()})
			return
		} else {
			c.JSON(200, gin.H{"data": "删除成功"})
		}
	}
}

func (r *ResourceCtl) Create() func(c *gin.Context) {
	fmt.Println("create")
	return func(c *gin.Context) {
		var resource = c.Param("resource")

		type ResouceParam struct {
			Yaml string `json:"yaml"`
		}

		var param ResouceParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(400, gin.H{"error": "解析请求体失败: " + err.Error()})
			return
		}

		err := r.resourceService.CreateResource(resource, param.Yaml)
		if err != nil {
			c.JSON(400, gin.H{"error": "创建失败：" + err.Error()})
			return
		} else {
			c.JSON(200, gin.H{"data": "创建成功"})
		}
	}
}

func (r *ResourceCtl) GetGVR() func(c *gin.Context) {
	return func(c *gin.Context) {
		var resource = c.Query("resource")

		gvr, err := r.resourceService.GetGVR(resource)
		if err != nil {
			c.JSON(400, gin.H{"error": "资源错误：" + err.Error()})
			return
		} else {
			c.JSON(200, gin.H{"data": *gvr})
		}
	}
}
