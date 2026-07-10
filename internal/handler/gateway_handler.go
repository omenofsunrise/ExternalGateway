package handler

import (
	"context"
	"log"

	"external-gateway/internal/domain"

	pb "external-gateway/api/generated/gateway"
)

type GatewayHandler struct {
	pb.UnimplementedExternalGatewayServer
	geminiService domain.GeminiService
}

func NewGatewayHandler(geminiService domain.GeminiService) *GatewayHandler {
	return &GatewayHandler{
		geminiService: geminiService,
	}
}

func (h *GatewayHandler) AskGemini(ctx context.Context, req *pb.AskGeminiRequest) (*pb.AskGeminiResponse, error) {
	log.Printf("Received AskGemini request with prompt: %s", req.Prompt)

	domainReq := &domain.AskGeminiRequest{
		Prompt: req.Prompt,
	}

	domainResp, err := h.geminiService.AskGemini(domainReq)
	if err != nil {
		log.Printf("Error processing AskGemini: %v", err)
		return &pb.AskGeminiResponse{
			Content: "",
			Success: false,
		}, nil
	}

	return &pb.AskGeminiResponse{
		Content: domainResp.Content,
		Success: domainResp.Success,
	}, nil
}
