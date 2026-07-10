package service

import (
	"context"
	"fmt"

	"external-gateway/internal/adapter/gemini"
	"external-gateway/internal/domain"
)

type GeminiServiceConfig struct {
	SystemInstruction string
	Temperature       float64
	MaxTokens         int32
}

type GeminiService struct {
	geminiClient *gemini.Client
	config       *GeminiServiceConfig
}

func NewGeminiService(client *gemini.Client, config *GeminiServiceConfig) *GeminiService {
	if config == nil {
		config = &GeminiServiceConfig{
			Temperature: 0.7,
			MaxTokens:   1000,
		}
	}
	return &GeminiService{
		geminiClient: client,
		config:       config,
	}
}

func (s *GeminiService) AskGemini(req *domain.AskGeminiRequest) (*domain.AskGeminiResponse, error) {
	if req.Prompt == "" {
		return &domain.AskGeminiResponse{
			Content: "",
			Success: false,
		}, fmt.Errorf("prompt cannot be empty")
	}

	opts := []gemini.RequestOption{
		gemini.WithGenerationConfig(gemini.GenerationConfig{
			Temperature:     &s.config.Temperature,
			MaxOutputTokens: &s.config.MaxTokens,
		}),
	}

	if s.config.SystemInstruction != "" {
		opts = append(opts, gemini.WithSystemInstruction(s.config.SystemInstruction))
	}

	resp, err := s.geminiClient.GenerateContent(
		context.Background(),
		req.Prompt,
		opts...,
	)
	if err != nil {
		return s.handleError(err)
	}

	text := s.geminiClient.GetTextFromResponse(resp)
	if text == "" {
		return s.handleEmptyResponse(resp)
	}

	tokens, _ := s.geminiClient.CountTokens(context.Background(), req.Prompt, opts...)
	fmt.Printf("Used tokens: %d\n", tokens)

	return &domain.AskGeminiResponse{
		Content: text,
		Success: true,
	}, nil
}

func (s *GeminiService) AskGeminiWithHistory(req *domain.AskGeminiRequest, history []domain.ChatMessage) (*domain.AskGeminiResponse, error) {
	if req.Prompt == "" {
		return &domain.AskGeminiResponse{
			Content: "",
			Success: false,
		}, fmt.Errorf("prompt cannot be empty")
	}

	var contents []gemini.Content
	for _, msg := range history {
		contents = append(contents, gemini.Content{
			Role: msg.Role,
			Parts: []gemini.Part{
				{Text: msg.Content},
			},
		})
	}

	opts := []gemini.RequestOption{
		gemini.WithGenerationConfig(gemini.GenerationConfig{
			Temperature:     &s.config.Temperature,
			MaxOutputTokens: &s.config.MaxTokens,
		}),
	}

	if s.config.SystemInstruction != "" {
		opts = append(opts, gemini.WithSystemInstruction(s.config.SystemInstruction))
	}

	resp, err := s.geminiClient.GenerateContentWithHistory(
		context.Background(),
		contents,
		[]gemini.Part{{Text: req.Prompt}},
		opts...,
	)
	if err != nil {
		return s.handleError(err)
	}

	text := s.geminiClient.GetTextFromResponse(resp)
	if text == "" {
		return s.handleEmptyResponse(resp)
	}

	return &domain.AskGeminiResponse{
		Content: text,
		Success: true,
	}, nil
}

func (s *GeminiService) AskGeminiStream(req *domain.AskGeminiRequest) (<-chan string, <-chan error, error) {
	if req.Prompt == "" {
		return nil, nil, fmt.Errorf("prompt cannot be empty")
	}

	opts := []gemini.RequestOption{
		gemini.WithGenerationConfig(gemini.GenerationConfig{
			Temperature:     &s.config.Temperature,
			MaxOutputTokens: &s.config.MaxTokens,
		}),
	}

	if s.config.SystemInstruction != "" {
		opts = append(opts, gemini.WithSystemInstruction(s.config.SystemInstruction))
	}

	stream, err := s.geminiClient.GenerateContentStream(
		context.Background(),
		req.Prompt,
		opts...,
	)
	if err != nil {
		return nil, nil, err
	}

	textChan := make(chan string)
	errChan := make(chan error)

	go func() {
		defer close(textChan)
		defer close(errChan)

		for {
			select {
			case event, ok := <-stream.Events():
				if !ok {
					return
				}
				for _, part := range event.Candidate.Content.Parts {
					if part.Text != "" {
						textChan <- part.Text
					}
				}
				if event.IsComplete {
					return
				}
			case err, ok := <-stream.Errors():
				if ok {
					errChan <- err
					return
				}
			}
		}
	}()

	return textChan, errChan, nil
}

func (s *GeminiService) handleError(err error) (*domain.AskGeminiResponse, error) {
	if geminiErr, ok := err.(*gemini.GeminiError); ok {
		return &domain.AskGeminiResponse{
			Content: fmt.Sprintf("Gemini API error: %s", geminiErr.Message),
			Success: false,
		}, err
	}
	return &domain.AskGeminiResponse{
		Content: fmt.Sprintf("Error: %v", err),
		Success: false,
	}, err
}

func (s *GeminiService) handleEmptyResponse(resp *gemini.GenerateContentResponse) (*domain.AskGeminiResponse, error) {
	if resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != "" {
		msg := fmt.Sprintf("Prompt blocked: %s", resp.PromptFeedback.BlockReason)
		return &domain.AskGeminiResponse{
			Content: msg,
			Success: false,
		}, fmt.Errorf(msg)
	}
	return &domain.AskGeminiResponse{
		Content: "Empty response from Gemini",
		Success: false,
	}, fmt.Errorf("empty response from Gemini")
}
