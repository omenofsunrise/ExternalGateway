package service

import (
	"context"
	"external-gateway/internal/adapter/pastebin"
	"external-gateway/internal/domain"
	"fmt"
)

type PastebinService struct {
	pastebinClient pastebin.Client
}

func NewPastebinService(client *pastebin.Client) *PastebinService {
	return &PastebinService{pastebinClient: *client}
}

func (s *PastebinService) AskTextPrompt(req domain.MakePromptRequest) (domain.MakePromptResponse, error) {
	if req.Prompt == "" {
		return domain.MakePromptResponse{Success: false, Content: ""}, fmt.Errorf("Prompt cannot be empty")
	}

	res, err := s.pastebinClient.SimpleChat(context.Background(), req.Prompt)
	if err != nil {
		return domain.MakePromptResponse{Success: false, Content: ""}, fmt.Errorf(err.Error())
	}

	return domain.MakePromptResponse{Success: true, Content: res}, nil
}
