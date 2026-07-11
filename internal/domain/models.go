package domain

import "external-gateway/internal/adapter/gemini"

type MakePromptRequest struct {
	Prompt string
}

type MakePromptResponse struct {
	Content string
	Success bool
}

type ChatMessage struct {
	Role    gemini.Role
	Content string
}
