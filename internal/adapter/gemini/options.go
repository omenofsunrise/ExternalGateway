package gemini

type RequestOption func(*GenerateContentRequest)

func WithSystemInstruction(instruction string) RequestOption {
	return func(req *GenerateContentRequest) {
		req.SystemInstruction = &Content{
			Parts: []Part{{Text: instruction}},
		}
	}
}

func WithGenerationConfig(config GenerationConfig) RequestOption {
	return func(req *GenerateContentRequest) {
		req.GenerationConfig = &config
	}
}

func WithTools(tools ...Tool) RequestOption {
	return func(req *GenerateContentRequest) {
		req.Tools = append(req.Tools, tools...)
	}
}

func WithToolConfig(mode FunctionCallingMode, allowedNames ...string) RequestOption {
	return func(req *GenerateContentRequest) {
		req.ToolConfig = &ToolConfig{
			FunctionCallingConfig: &FunctionCallingConfig{
				Mode:         mode,
				AllowedNames: allowedNames,
			},
		}
	}
}

func WithSafetySettings(settings ...SafetySetting) RequestOption {
	return func(req *GenerateContentRequest) {
		req.SafetySettings = append(req.SafetySettings, settings...)
	}
}

func WithTemperature(temp float64) RequestOption {
	return func(req *GenerateContentRequest) {
		if req.GenerationConfig == nil {
			req.GenerationConfig = &GenerationConfig{}
		}
		req.GenerationConfig.Temperature = &temp
	}
}

func WithMaxTokens(maxTokens int32) RequestOption {
	return func(req *GenerateContentRequest) {
		if req.GenerationConfig == nil {
			req.GenerationConfig = &GenerationConfig{}
		}
		req.GenerationConfig.MaxOutputTokens = &maxTokens
	}
}

func WithTopP(topP float64) RequestOption {
	return func(req *GenerateContentRequest) {
		if req.GenerationConfig == nil {
			req.GenerationConfig = &GenerationConfig{}
		}
		req.GenerationConfig.TopP = &topP
	}
}

func WithTopK(topK float64) RequestOption {
	return func(req *GenerateContentRequest) {
		if req.GenerationConfig == nil {
			req.GenerationConfig = &GenerationConfig{}
		}
		req.GenerationConfig.TopK = &topK
	}
}

func WithStopSequences(sequences ...string) RequestOption {
	return func(req *GenerateContentRequest) {
		if req.GenerationConfig == nil {
			req.GenerationConfig = &GenerationConfig{}
		}
		req.GenerationConfig.StopSequences = sequences
	}
}

func WithCandidateCount(count int32) RequestOption {
	return func(req *GenerateContentRequest) {
		if req.GenerationConfig == nil {
			req.GenerationConfig = &GenerationConfig{}
		}
		req.GenerationConfig.CandidateCount = &count
	}
}

func WithBlockNone() RequestOption {
	return WithSafetySettings(
		SafetySetting{Category: SafetyCategoryHarassment, Threshold: SafetyThresholdBlockNone},
		SafetySetting{Category: SafetyCategoryHateSpeech, Threshold: SafetyThresholdBlockNone},
		SafetySetting{Category: SafetyCategorySexuallyExplicit, Threshold: SafetyThresholdBlockNone},
		SafetySetting{Category: SafetyCategoryDangerousContent, Threshold: SafetyThresholdBlockNone},
	)
}

func WithBlockHigh() RequestOption {
	return WithSafetySettings(
		SafetySetting{Category: SafetyCategoryHarassment, Threshold: SafetyThresholdBlockHigh},
		SafetySetting{Category: SafetyCategoryHateSpeech, Threshold: SafetyThresholdBlockHigh},
		SafetySetting{Category: SafetyCategorySexuallyExplicit, Threshold: SafetyThresholdBlockHigh},
		SafetySetting{Category: SafetyCategoryDangerousContent, Threshold: SafetyThresholdBlockHigh},
	)
}

func WithDefaultSafety() RequestOption {
	return WithSafetySettings(
		SafetySetting{Category: SafetyCategoryHarassment, Threshold: SafetyThresholdBlockMedium},
		SafetySetting{Category: SafetyCategoryHateSpeech, Threshold: SafetyThresholdBlockMedium},
		SafetySetting{Category: SafetyCategorySexuallyExplicit, Threshold: SafetyThresholdBlockMedium},
		SafetySetting{Category: SafetyCategoryDangerousContent, Threshold: SafetyThresholdBlockMedium},
	)
}

func NewFunctionTool(name, description string, parameters *Parameters) Tool {
	return Tool{
		FunctionDeclarations: []FunctionDeclaration{
			{
				Name:        name,
				Description: description,
				Parameters:  parameters,
			},
		},
	}
}

func NewFunctionToolWithMultiple(functions ...FunctionDeclaration) Tool {
	return Tool{
		FunctionDeclarations: functions,
	}
}

func NewParameters(properties map[string]Property, required ...string) *Parameters {
	return &Parameters{
		Type:       "OBJECT",
		Properties: properties,
		Required:   required,
	}
}

func NewStringProperty(description string, enum ...string) Property {
	prop := Property{
		Type:        "string",
		Description: description,
	}
	if len(enum) > 0 {
		prop.Enum = enum
	}
	return prop
}

func NewNumberProperty(description string) Property {
	return Property{
		Type:        "number",
		Description: description,
	}
}

func NewIntegerProperty(description string) Property {
	return Property{
		Type:        "integer",
		Description: description,
	}
}

func NewBooleanProperty(description string) Property {
	return Property{
		Type:        "boolean",
		Description: description,
	}
}

func NewArrayProperty(itemsType string, description string) Property {
	return Property{
		Type:        "array",
		Description: description,
		Items: &Property{
			Type: itemsType,
		},
	}
}

func NewImagePart(data []byte, mimeType string) Part {
	return Part{
		InlineData: &InlineData{
			MimeType: mimeType,
			Data:     string(data),
		},
	}
}

func NewImagePartFromBase64(data string, mimeType string) Part {
	return Part{
		InlineData: &InlineData{
			MimeType: mimeType,
			Data:     data,
		},
	}
}

func NewTextPart(text string) Part {
	return Part{
		Text: text,
	}
}

func NewFunctionCallPart(name string, args map[string]interface{}) Part {
	return Part{
		FunctionCall: &FunctionCall{
			Name: name,
			Args: args,
		},
	}
}

func NewFunctionResponsePart(name string, response map[string]interface{}) Part {
	return Part{
		FunctionResponse: &FunctionResponse{
			Name:     name,
			Response: response,
		},
	}
}

func NewContent(role Role, parts ...Part) Content {
	return Content{
		Role:  role,
		Parts: parts,
	}
}

func NewUserContent(parts ...Part) Content {
	return Content{
		Role:  RoleUser,
		Parts: parts,
	}
}

func NewUserTextContent(text string) Content {
	return Content{
		Role:  RoleUser,
		Parts: []Part{{Text: text}},
	}
}

func NewModelContent(parts ...Part) Content {
	return Content{
		Role:  RoleModel,
		Parts: parts,
	}
}

func NewModelTextContent(text string) Content {
	return Content{
		Role:  RoleModel,
		Parts: []Part{{Text: text}},
	}
}

func DefaultGenerationConfig() GenerationConfig {
	temp := 0.7
	maxTokens := int32(1000)
	return GenerationConfig{
		Temperature:     &temp,
		MaxOutputTokens: &maxTokens,
	}
}

func CreativeGenerationConfig() GenerationConfig {
	temp := 0.9
	maxTokens := int32(2000)
	topP := 0.95
	return GenerationConfig{
		Temperature:     &temp,
		MaxOutputTokens: &maxTokens,
		TopP:            &topP,
	}
}

func PreciseGenerationConfig() GenerationConfig {
	temp := 0.1
	maxTokens := int32(500)
	topP := 0.1
	return GenerationConfig{
		Temperature:     &temp,
		MaxOutputTokens: &maxTokens,
		TopP:            &topP,
	}
}
