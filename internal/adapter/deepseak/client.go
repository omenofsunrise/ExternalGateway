package deepseak

import (
	"encoding/json"
	"external-gateway/internal/adapter/common"
	"fmt"
	"net/http"
)

type DeepseakClient struct {
	common.BaseAIClient
}

type BalanceInfo struct {
	Currency        string `json:"currency"`
	TotalBalance    string `json:"total_balance"`
	GrantedBalance  string `json:"granted_balance"`
	ToppedUpBalance string `json:"topped_up_balance"`
}

type BalanceResponse struct {
	IsAvailable  bool          `json:"is_available"`
	BalanceInfos []BalanceInfo `json:"balance_infos"`
}

func NewClient(httpClient *http.Client, apiKey string) *DeepseakClient {
	baseUrl := "https://api.deepseek.com"
	model := "deepseek-v4-flash"
	client := common.NewBaseAIClient(httpClient, apiKey, baseUrl, model)
	return &DeepseakClient{BaseAIClient: *client}
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
