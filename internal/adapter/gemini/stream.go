package gemini

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StreamEvent struct {
	Candidate    Candidate
	Index        int32
	IsComplete   bool
	FinishReason string
}

type Stream struct {
	client  *Client
	ctx     context.Context
	events  chan StreamEvent
	errChan chan error
}

func (c *Client) GenerateContentStream(ctx context.Context, prompt string, opts ...RequestOption) (*Stream, error) {
	return c.GenerateContentStreamWithParts(ctx, []Part{{Text: prompt}}, opts...)
}

func (c *Client) GenerateContentStreamWithParts(ctx context.Context, parts []Part, opts ...RequestOption) (*Stream, error) {
	req := &GenerateContentRequest{
		Contents: []Content{
			{
				Role:  "user",
				Parts: parts,
			},
		},
	}

	for _, opt := range opts {
		opt(req)
	}

	stream := &Stream{
		client:  c,
		ctx:     ctx,
		events:  make(chan StreamEvent, 10),
		errChan: make(chan error, 1),
	}

	go stream.start(req)
	return stream, nil
}

func (s *Stream) start(req *GenerateContentRequest) {
	defer close(s.events)
	defer close(s.errChan)

	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?key=%s", s.client.baseURL, s.client.model, s.client.apiKey)

	jsonData, err := json.Marshal(req)
	if err != nil {
		s.errChan <- fmt.Errorf("failed to marshal request: %w", err)
		return
	}

	httpReq, err := http.NewRequestWithContext(s.ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.errChan <- fmt.Errorf("failed to create request: %w", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		s.errChan <- fmt.Errorf("failed to call Gemini API: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.errChan <- &GeminiError{
			StatusCode: resp.StatusCode,
			Message:    "stream request failed",
			Body:       string(body),
		}
		return
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			s.errChan <- fmt.Errorf("failed to read stream: %w", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		line = strings.TrimPrefix(line, "data: ")

		var streamResp struct {
			Candidates []Candidate `json:"candidates"`
		}

		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			s.errChan <- fmt.Errorf("failed to parse stream data: %w", err)
			return
		}

		for _, candidate := range streamResp.Candidates {
			event := StreamEvent{
				Candidate:    candidate,
				Index:        candidate.Index,
				IsComplete:   candidate.FinishReason != "",
				FinishReason: candidate.FinishReason.String(),
			}
			s.events <- event
		}
	}
}

func (s *Stream) Events() <-chan StreamEvent {
	return s.events
}

func (s *Stream) Errors() <-chan error {
	return s.errChan
}

func (c *Client) CollectStreamText(stream *Stream) (string, error) {
	var fullText string

	for {
		select {
		case event, ok := <-stream.Events():
			if !ok {
				return fullText, nil
			}
			for _, part := range event.Candidate.Content.Parts {
				fullText += part.Text
			}
			if event.IsComplete {
				return fullText, nil
			}
		case err, ok := <-stream.Errors():
			if !ok {
				return fullText, nil
			}
			return fullText, err
		case <-stream.ctx.Done():
			return fullText, stream.ctx.Err()
		}
	}
}
