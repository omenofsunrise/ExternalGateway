package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"external-gateway/internal/adapter/common"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GeminiClient struct {
	*common.BaseAIClient
}

func NewClient(apiKey, model string, opts ...ClientOption) *GeminiClient {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	client := common.NewBaseAIClient(httpClient, apiKey, "https://generativelanguage.googleapis.com/v1beta", model)

	for _, opt := range opts {
		opt(client)
	}

	return &GeminiClient{BaseAIClient: client}
}

type ClientOption func(*common.BaseAIClient)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *common.BaseAIClient) {
		c.BaseClient.HttpClient = httpClient
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *common.BaseAIClient) {
		c.BaseClient.BaseUrl = baseURL
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *common.BaseAIClient) {
		c.BaseClient.HttpClient.Timeout = timeout
	}
}

func (c *GeminiClient) GenerateContent(ctx context.Context, prompt string, opts ...RequestOption) (*GenerateContentResponse, error) {
	return c.GenerateContentWithParts(ctx, []Part{{Text: prompt}}, opts...)
}

func (c *GeminiClient) GenerateContentWithParts(ctx context.Context, parts []Part, opts ...RequestOption) (*GenerateContentResponse, error) {
	req := &GenerateContentRequest{
		Contents: []Content{
			{
				Role:  RoleUser,
				Parts: parts,
			},
		},
	}

	for _, opt := range opts {
		opt(req)
	}

	return c.generateContent(ctx, req)
}

func (c *GeminiClient) GenerateContentWithHistory(ctx context.Context, history []Content, newParts []Part, opts ...RequestOption) (*GenerateContentResponse, error) {
	contents := make([]Content, len(history)+1)
	copy(contents, history)
	contents[len(history)] = Content{
		Role:  RoleUser,
		Parts: newParts,
	}

	req := &GenerateContentRequest{
		Contents: contents,
	}

	for _, opt := range opts {
		opt(req)
	}

	return c.generateContent(ctx, req)
}

func (c *GeminiClient) generateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error) {
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.BaseClient.BaseUrl, c.Model, c.BaseClient)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &GeminiError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("API request failed with status: %s", resp.Status),
			Body:       string(body),
		}
	}

	var geminiResp GenerateContentResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		if geminiResp.PromptFeedback != nil && geminiResp.PromptFeedback.BlockReason != "" {
			return nil, fmt.Errorf("prompt blocked: %s", geminiResp.PromptFeedback.BlockReason)
		}
		return nil, fmt.Errorf("empty response from Gemini")
	}

	return &geminiResp, nil
}

func (c *GeminiClient) GenerateContentWithFunctions(
	ctx context.Context,
	prompt string,
	functions []Tool,
	opts ...RequestOption,
) (*GenerateContentResponse, error) {
	allOpts := append([]RequestOption{WithTools(functions...)}, opts...)
	return c.GenerateContent(ctx, prompt, allOpts...)
}

func (c *GeminiClient) HandleFunctionCall(
	originalResponse *GenerateContentResponse,
	functionName string,
	result map[string]interface{},
) (*GenerateContentRequest, error) {
	if len(originalResponse.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	candidate := originalResponse.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("no parts in candidate")
	}

	var functionCall *FunctionCall
	for _, part := range candidate.Content.Parts {
		if part.FunctionCall != nil && part.FunctionCall.Name == functionName {
			functionCall = part.FunctionCall
			break
		}
	}

	if functionCall == nil {
		return nil, fmt.Errorf("function %s not found in response", functionName)
	}

	return &GenerateContentRequest{
		Contents: []Content{
			{
				Role:  RoleUser,
				Parts: []Part{{Text: fmt.Sprintf("Call function: %s", functionName)}},
			},
			{
				Role:  RoleModel,
				Parts: []Part{{FunctionCall: functionCall}},
			},
			{
				Role: RoleUser,
				Parts: []Part{
					{
						FunctionResponse: &FunctionResponse{
							Name:     functionName,
							Response: result,
						},
					},
				},
			},
		},
	}, nil
}

func (c *GeminiClient) HandleFunctionCallWithHistory(
	history []Content,
	originalResponse *GenerateContentResponse,
	functionName string,
	result map[string]interface{},
) (*GenerateContentRequest, error) {
	if len(originalResponse.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	candidate := originalResponse.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("no parts in candidate")
	}

	var functionCall *FunctionCall
	for _, part := range candidate.Content.Parts {
		if part.FunctionCall != nil && part.FunctionCall.Name == functionName {
			functionCall = part.FunctionCall
			break
		}
	}

	if functionCall == nil {
		return nil, fmt.Errorf("function %s not found in response", functionName)
	}

	contents := make([]Content, len(history)+2)
	copy(contents, history)
	contents[len(history)] = Content{
		Role:  RoleModel,
		Parts: []Part{{FunctionCall: functionCall}},
	}
	contents[len(history)+1] = Content{
		Role: RoleUser,
		Parts: []Part{
			{
				FunctionResponse: &FunctionResponse{
					Name:     functionName,
					Response: result,
				},
			},
		},
	}

	return &GenerateContentRequest{
		Contents: contents,
	}, nil
}

func (c *GeminiClient) CountTokens(ctx context.Context, prompt string, opts ...RequestOption) (int32, error) {
	req := &CountTokensRequest{
		Contents: []Content{
			{
				Role:  RoleUser,
				Parts: []Part{{Text: prompt}},
			},
		},
	}

	tempReq := &GenerateContentRequest{}
	for _, opt := range opts {
		opt(tempReq)
	}
	if tempReq.SystemInstruction != nil {
		req.SystemInstruction = tempReq.SystemInstruction
	}
	if len(tempReq.Tools) > 0 {
		req.Tools = tempReq.Tools
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:countTokens?key=%s", c.BaseUrl, c.Model, c.ApiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, &GeminiError{
			StatusCode: resp.StatusCode,
			Message:    "failed to count tokens",
			Body:       string(body),
		}
	}

	var tokenResp CountTokensResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return tokenResp.TotalTokens, nil
}

func (c *GeminiClient) CountTokensWithContent(ctx context.Context, content Content, opts ...RequestOption) (int32, error) {
	req := &CountTokensRequest{
		Contents: []Content{content},
	}

	tempReq := &GenerateContentRequest{}
	for _, opt := range opts {
		opt(tempReq)
	}
	if tempReq.SystemInstruction != nil {
		req.SystemInstruction = tempReq.SystemInstruction
	}
	if len(tempReq.Tools) > 0 {
		req.Tools = tempReq.Tools
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:countTokens?key=%s", c.BaseUrl, c.Model, c.ApiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, &GeminiError{
			StatusCode: resp.StatusCode,
			Message:    "failed to count tokens",
			Body:       string(body),
		}
	}

	var tokenResp CountTokensResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return tokenResp.TotalTokens, nil
}

func (c *GeminiClient) EmbedContent(ctx context.Context, text string) ([]float32, error) {
	return c.EmbedContentWithConfig(ctx, Content{
		Parts: []Part{{Text: text}},
	}, "")
}

func (c *GeminiClient) EmbedContentWithConfig(
	ctx context.Context,
	content Content,
	taskType EmbeddingTaskType,
) ([]float32, error) {
	req := &EmbedContentRequest{
		Model:    c.Model,
		Content:  content,
		TaskType: taskType,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:embedContent?key=%s", c.BaseUrl, c.Model, c.ApiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &GeminiError{
			StatusCode: resp.StatusCode,
			Message:    "failed to create embedding",
			Body:       string(body),
		}
	}

	var embedResp EmbedContentResponse
	if err := json.Unmarshal(body, &embedResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if embedResp.Embedding == nil {
		return nil, fmt.Errorf("empty embedding response")
	}

	return embedResp.Embedding.Values, nil
}

func (c *GeminiClient) GetTextFromResponse(resp *GenerateContentResponse) string {
	if len(resp.Candidates) == 0 {
		return ""
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			return part.Text
		}
	}
	return ""
}

func (c *GeminiClient) GetFunctionCallsFromResponse(resp *GenerateContentResponse) []FunctionCall {
	var calls []FunctionCall
	if len(resp.Candidates) == 0 {
		return calls
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.FunctionCall != nil {
			calls = append(calls, *part.FunctionCall)
		}
	}
	return calls
}

func (c *GeminiClient) GetAllTextFromResponse(resp *GenerateContentResponse) []string {
	var texts []string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.Text != "" {
				texts = append(texts, part.Text)
			}
		}
	}
	return texts
}

func (c *GeminiClient) GetFinishReason(resp *GenerateContentResponse) (FinishReason, error) {
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}
	return resp.Candidates[0].FinishReason, nil
}

func (c *GeminiClient) IsResponseBlocked(resp *GenerateContentResponse) bool {
	if resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != "" {
		return true
	}
	if len(resp.Candidates) > 0 && resp.Candidates[0].FinishReason.IsError() {
		return true
	}
	return false
}

func (c *GeminiClient) BatchGenerateContent(ctx context.Context, prompts []string, opts ...RequestOption) ([]*GenerateContentResponse, error) {
	responses := make([]*GenerateContentResponse, len(prompts))
	for i, prompt := range prompts {
		resp, err := c.GenerateContent(ctx, prompt, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to generate content for prompt %d: %w", i, err)
		}
		responses[i] = resp
	}
	return responses, nil
}

func (c *GeminiClient) GenerateContentWithRetry(
	ctx context.Context,
	prompt string,
	maxRetries int,
	opts ...RequestOption,
) (*GenerateContentResponse, error) {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := c.GenerateContent(ctx, prompt, opts...)
		if err == nil {
			return resp, nil
		}

		if geminiErr, ok := err.(*GeminiError); ok && geminiErr.IsRetryable() {
			lastErr = err
			delay := time.Duration(1<<attempt) * time.Second
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			time.Sleep(delay)
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
