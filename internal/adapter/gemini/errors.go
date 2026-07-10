package gemini

import "fmt"

type GeminiError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *GeminiError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("Gemini API error (status %d): %s - %s", e.StatusCode, e.Message, e.Body)
	}
	return fmt.Sprintf("Gemini API error (status %d): %s", e.StatusCode, e.Message)
}

func (e *GeminiError) IsRetryable() bool {
	return e.StatusCode >= 500 || e.StatusCode == 429
}

func (e *GeminiError) IsAuthError() bool {
	return e.StatusCode == 401 || e.StatusCode == 403
}
