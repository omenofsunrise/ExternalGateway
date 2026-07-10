package domain

import "external-gateway/internal/adapter/gemini"

type AskGeminiRequest struct {
	Prompt string
}

type AskGeminiResponse struct {
	Content string
	Success bool
}

type ChatMessage struct {
	Role    gemini.Role
	Content string
}
