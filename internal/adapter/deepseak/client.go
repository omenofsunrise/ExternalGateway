package deepseak

import (
	"bytes"
	"encoding/json"
	"external-gateway/internal/adapter/common"
	"fmt"
	"net/http"
)

func NewClient(httpClient *http.Client, apiKey string) *DeepseakClient {
	baseUrl := "https://api.deepseek.com"
	model := "deepseek-v4-flash"
	client := common.NewBaseAIClient(httpClient, apiKey, baseUrl, model)
	return &DeepseakClient{BaseAIClient: *client}
}

func (c *DeepseakClient) CreateChatCompletion(req *CompletionRequest) (*ChatCompletion, error) {
	url := fmt.Sprintf("%s/chat/completions", c.BaseUrl)
	if req.Model == "" {
		req.Model = c.Model
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error serializing body to json: %s", err)
	}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error preparing http request: %s", err)
	}

	common.AddAuthHeader(httpReq, c.ApiKey)
	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var completion *ChatCompletion
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	return completion, nil
}

func (c *DeepseakClient) GetBalance() (*BalanceResponse, error) {
	url := fmt.Sprintf("%v/user/balance", c.BaseUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	common.AddAuthHeader(req, c.ApiKey)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make http request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %v", resp.StatusCode)
	}

	var balanceResponse BalanceResponse
	if err = json.NewDecoder(resp.Body).Decode(&balanceResponse); err != nil {
		return nil, fmt.Errorf("error serialize json: %v", err)
	}
	return &balanceResponse, nil
}
