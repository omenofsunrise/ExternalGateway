package domain

type GeminiService interface {
	AskGemini(req *MakePromptRequest) (*MakePromptResponse, error)
}
