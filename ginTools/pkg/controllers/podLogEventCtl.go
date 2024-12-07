package controllers

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xingyunyang01/ginTools/pkg/services"
)

type PodLogEventCtl struct {
	podLogEventService *services.PodLogEventService
}

func NewPodLogEventCtl(service *services.PodLogEventService) *PodLogEventCtl {
	return &PodLogEventCtl{podLogEventService: service}
}

func (p *PodLogEventCtl) GetLog() func(c *gin.Context) {
	return func(c *gin.Context) {
		ns := c.DefaultQuery("ns", "default")
		podname := c.DefaultQuery("podname", "")

		var tailLine int64 = 100

		req := p.podLogEventService.GetLogs(ns, podname, tailLine)

		rc, err := req.Stream(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		defer rc.Close()

		logData, err := ioutil.ReadAll(rc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"data": string(logData)})
	}
}

func (p *PodLogEventCtl) GetEvent() func(c *gin.Context) {
	return func(c *gin.Context) {
		ns := c.DefaultQuery("ns", "default")
		podname := c.DefaultQuery("podname", "")

		e, err := p.podLogEventService.GetEvents(ns, podname)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"data": e})
	}
}
