package gemini

type Role string

const (
	RoleUser  Role = "user"
	RoleModel Role = "model"
)

type FinishReason string

const (
	FinishReasonStop              FinishReason = "STOP"
	FinishReasonMaxTokens         FinishReason = "MAX_TOKENS"
	FinishReasonSafety            FinishReason = "SAFETY"
	FinishReasonRecitation        FinishReason = "RECITATION"
	FinishReasonOther             FinishReason = "OTHER"
	FinishReasonBlocklist         FinishReason = "BLOCKLIST"
	FinishReasonProhibitedContent FinishReason = "PROHIBITED_CONTENT"
	FinishReasonSpii              FinishReason = "SPII"
)

type FunctionCallingMode string

const (
	FunctionCallingModeNone FunctionCallingMode = "NONE"
	FunctionCallingModeAuto FunctionCallingMode = "AUTO"
	FunctionCallingModeAny  FunctionCallingMode = "ANY"
)

type SafetyCategory string

const (
	SafetyCategoryHarassment       SafetyCategory = "HARM_CATEGORY_HARASSMENT"
	SafetyCategoryHateSpeech       SafetyCategory = "HARM_CATEGORY_HATE_SPEECH"
	SafetyCategorySexuallyExplicit SafetyCategory = "HARM_CATEGORY_SEXUALLY_EXPLICIT"
	SafetyCategoryDangerousContent SafetyCategory = "HARM_CATEGORY_DANGEROUS_CONTENT"
	SafetyCategoryCivicIntegrity   SafetyCategory = "HARM_CATEGORY_CIVIC_INTEGRITY"
)

type SafetyThreshold string

const (
	SafetyThresholdBlockNone   SafetyThreshold = "BLOCK_NONE"
	SafetyThresholdBlockLow    SafetyThreshold = "BLOCK_LOW_AND_ABOVE"
	SafetyThresholdBlockMedium SafetyThreshold = "BLOCK_MEDIUM_AND_ABOVE"
	SafetyThresholdBlockHigh   SafetyThreshold = "BLOCK_HIGH_AND_ABOVE"
)

type SafetyProbability string

const (
	SafetyProbabilityNegligible SafetyProbability = "NEGLIGIBLE"
	SafetyProbabilityLow        SafetyProbability = "LOW"
	SafetyProbabilityMedium     SafetyProbability = "MEDIUM"
	SafetyProbabilityHigh       SafetyProbability = "HIGH"
)

type EmbeddingTaskType string

const (
	EmbeddingTaskTypeRetrievalQuery     EmbeddingTaskType = "RETRIEVAL_QUERY"
	EmbeddingTaskTypeRetrievalDocument  EmbeddingTaskType = "RETRIEVAL_DOCUMENT"
	EmbeddingTaskTypeSemanticSimilarity EmbeddingTaskType = "SEMANTIC_SIMILARITY"
	EmbeddingTaskTypeClassification     EmbeddingTaskType = "CLASSIFICATION"
	EmbeddingTaskTypeClustering         EmbeddingTaskType = "CLUSTERING"
	EmbeddingTaskTypeQuestionAnswering  EmbeddingTaskType = "QUESTION_ANSWERING"
	EmbeddingTaskTypeFactVerification   EmbeddingTaskType = "FACT_VERIFICATION"
)

type Content struct {
	Role  Role   `json:"role,omitempty"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text             string            `json:"text,omitempty"`
	InlineData       *InlineData       `json:"inlineData,omitempty"`
	FunctionCall     *FunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *FunctionResponse `json:"functionResponse,omitempty"`
}

type InlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type GenerationConfig struct {
	Temperature     *float64 `json:"temperature,omitempty"`
	MaxOutputTokens *int32   `json:"maxOutputTokens,omitempty"`
	TopP            *float64 `json:"topP,omitempty"`
	TopK            *float64 `json:"topK,omitempty"`
	StopSequences   []string `json:"stopSequences,omitempty"`
	CandidateCount  *int32   `json:"candidateCount,omitempty"`
}

type Tool struct {
	FunctionDeclarations []FunctionDeclaration `json:"functionDeclarations"`
}

type FunctionDeclaration struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  *Parameters `json:"parameters,omitempty"`
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	Enum        []string  `json:"enum,omitempty"`
	Items       *Property `json:"items,omitempty"`
}

type FunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type FunctionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

type SafetySetting struct {
	Category  SafetyCategory  `json:"category"`
	Threshold SafetyThreshold `json:"threshold"`
}

type GenerateContentRequest struct {
	Contents          []Content         `json:"contents"`
	SystemInstruction *Content          `json:"systemInstruction,omitempty"`
	GenerationConfig  *GenerationConfig `json:"generationConfig,omitempty"`
	Tools             []Tool            `json:"tools,omitempty"`
	ToolConfig        *ToolConfig       `json:"toolConfig,omitempty"`
	SafetySettings    []SafetySetting   `json:"safetySettings,omitempty"`
}

type ToolConfig struct {
	FunctionCallingConfig *FunctionCallingConfig `json:"functionCallingConfig,omitempty"`
}

type FunctionCallingConfig struct {
	Mode         FunctionCallingMode `json:"mode,omitempty"`
	AllowedNames []string            `json:"allowedFunctionNames,omitempty"`
}

type GenerateContentResponse struct {
	Candidates     []Candidate     `json:"candidates"`
	PromptFeedback *PromptFeedback `json:"promptFeedback,omitempty"`
}

type Candidate struct {
	Content       Content        `json:"content"`
	FinishReason  FinishReason   `json:"finishReason,omitempty"`
	Index         int32          `json:"index,omitempty"`
	SafetyRatings []SafetyRating `json:"safetyRatings,omitempty"`
}

type SafetyRating struct {
	Category    SafetyCategory    `json:"category"`
	Probability SafetyProbability `json:"probability"`
}

type PromptFeedback struct {
	BlockReason   string         `json:"blockReason,omitempty"`
	SafetyRatings []SafetyRating `json:"safetyRatings,omitempty"`
}

type CountTokensRequest struct {
	Contents          []Content `json:"contents"`
	SystemInstruction *Content  `json:"systemInstruction,omitempty"`
	Tools             []Tool    `json:"tools,omitempty"`
}

type CountTokensResponse struct {
	TotalTokens int32 `json:"totalTokens"`
}

type EmbedContentRequest struct {
	Model    string            `json:"-"`
	Content  Content           `json:"content"`
	TaskType EmbeddingTaskType `json:"taskType,omitempty"`
	Title    string            `json:"title,omitempty"`
}

type EmbedContentResponse struct {
	Embedding *Embedding `json:"embedding"`
}

type Embedding struct {
	Values []float32 `json:"values"`
}
