package services

import (
	"github.com/xingyunyang01/APIAgent/pkg/core/agent"
	"github.com/xingyunyang01/APIAgent/pkg/models"
)

type ChatCompletionService struct {
	sc    *models.Config
	tools []models.ApiToolBundle
}

func NewChatCompletionService(sc *models.Config, tools []models.ApiToolBundle) *ChatCompletionService {
	return &ChatCompletionService{sc: sc, tools: tools}
}

func (s *ChatCompletionService) ChatCompletion(query string) (string, error) {
	return agent.Run(s.sc, s.tools, query)
}
