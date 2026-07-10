package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"external-gateway/internal/adapter/gemini"
	"external-gateway/internal/config"
	"external-gateway/internal/handler"
	"external-gateway/internal/server"
	"external-gateway/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Gemini.APIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	geminiClient := gemini.NewClient(
		cfg.Gemini.APIKey,
		cfg.Gemini.Model,
		gemini.WithTimeout(time.Duration(cfg.Gemini.Timeout)*time.Second),
	)

	geminiService := service.NewGeminiService(
		geminiClient,
		&service.GeminiServiceConfig{
			SystemInstruction: "Ты - полезный ассистент. Отвечай на русском языке, кратко и по делу.",
			Temperature:       0.7,
			MaxTokens:         1000,
		},
	)

	gatewayHandler := handler.NewGatewayHandler(geminiService)

	grpcServer := server.NewGRPCServer(cfg.Server.GRPCPort, gatewayHandler)

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")
	grpcServer.Stop()
}
