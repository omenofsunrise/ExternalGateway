package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server ServerConfig
	Gemini GeminiConfig
	Logger LoggerConfig
}

type ServerConfig struct {
	GRPCPort    int
	HTTPPort    int
	GatewayPort int
}

type GeminiConfig struct {
	APIKey  string
	Model   string
	Timeout int
}

type LoggerConfig struct {
	Level string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			GRPCPort: getEnvAsInt("GRPC_PORT", 50051),
		},
		Gemini: GeminiConfig{
			APIKey:  getEnv("GEMINI_API_KEY", ""),
			Model:   getEnv("GEMINI_MODEL", "gemini-2.0-flash-exp"),
			Timeout: getEnvAsInt("GEMINI_TIMEOUT", 30),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
