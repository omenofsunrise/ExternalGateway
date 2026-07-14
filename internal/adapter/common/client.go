package common

import "net/http"

type BaseClient struct {
	HttpClient *http.Client
	ApiKey     string
	BaseUrl    string
}

type BaseAIClient struct {
	BaseClient
	Model string
}

func NewBaseClient(httpClient *http.Client, apiKey, baseUrl string) *BaseClient {
	return &BaseClient{HttpClient: httpClient, ApiKey: apiKey, BaseUrl: baseUrl}
}

func NewBaseAIClient(httpClient *http.Client, apiKey, baseUrl, model string) *BaseAIClient {
	baseClient := NewBaseClient(httpClient, apiKey, baseUrl)
	return &BaseAIClient{BaseClient: *baseClient, Model: model}
}
