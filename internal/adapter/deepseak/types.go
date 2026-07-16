package deepseak

import "external-gateway/internal/adapter/common"

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

type CompletionRequest struct {
	Messages       []Message       `json:"messages"`
	Model          string          `json:"model"`
	MaxTokens      *int            `json:"max_tokens,omitempty"`
	Stream         *bool           `json:"stream,omitempty"`
	Temperature    *float32        `json:"temperature,omitempty"`     // from 0 to 2
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"` // json || text
	ToolChoice     *string         `json:"tool_choice,omitempty"`
}

type Message struct {
	Content    string  `json:"content"`
	Role       Role    `json:"role"`
	Name       *string `json:"name,omitempty"`
	ToolCallId *string `json:"tool_call_id,omitempty"`
}

type Tool struct {
	Type     string   `json:"type"` // "function"
	Function Function `json:"function"`
}

type Function struct {
	Name        string                 `json:"name"`
	Parameters  map[string]interface{} `json:"parameters"`
	Description *string                `json:"description,omitempty"`
	Strict      *bool                  `json:"strict,omitempty"` // beta
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type ChatCompletion struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

type ResponseMessage struct {
	Content   string     `json:"content"`
	Role      string     `json:"role"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Role string

const (
	SystemRole    Role = "system"
	UserRole      Role = "user"
	AssistantRole Role = "assistant"
	ToolMessage   Role = "tool"
)
