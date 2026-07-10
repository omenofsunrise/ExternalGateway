package domain

type GeminiService interface {
	AskGemini(req *AskGeminiRequest) (*AskGeminiResponse, error)
}
