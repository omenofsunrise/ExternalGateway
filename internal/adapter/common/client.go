package common

import (
	"fmt"
	"net/http"
)

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

func AddHeader(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

func AddAuthHeader(req *http.Request, apiKey string) {
	AddHeader(req)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
}
