package agent

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

// LLMClient handles communication with LLM providers
type LLMClient struct {
	apiKey   string
	baseURL  string
	model    string
	provider string
	client   *http.Client
}

// NewLLMClient creates a new LLM client
func NewLLMClient(apiKey, model, baseURL, provider string) *LLMClient {
	// Set default model based on provider
	if model == "" {
		if provider == "deepseek" {
			model = "deepseek-chat"
		} else {
			model = "gpt-4o-mini"
		}
	}

	// Set default base URL based on provider
	if baseURL == "" {
		if provider == "deepseek" {
			baseURL = "https://api.deepseek.com/v1"
		} else {
			baseURL = "https://api.openai.com/v1"
		}
	}

	return &LLMClient{
		apiKey:   apiKey,
		baseURL:  baseURL,
		model:    model,
		provider: provider,
		client:   &http.Client{},
	}
}

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

// ToolCall represents a tool call from the LLM
type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents the function to call
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool represents a tool definition
type Tool struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

// FunctionDef defines a function
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// LLMChatRequest represents the request to OpenAI
type LLMChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Tools       []Tool        `json:"tools,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// LLMChatResponse represents the response from OpenAI
type LLMChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// LLMChatStreamResponse represents a chunk from the streaming response
type LLMChatStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int                `json:"index"`
		Delta        LLMChatStreamDelta `json:"delta"`
		FinishReason *string            `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// LLMChatStreamDelta represents the delta content in a streaming response
type LLMChatStreamDelta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Chat sends a chat request to the LLM
func (c *LLMClient) Chat(messages []ChatMessage, tools []Tool) (*LLMChatResponse, error) {
	req := LLMChatRequest{
		Model:       c.model,
		Messages:    messages,
		Tools:       tools,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp LLMChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("LLM API error: %s (%s)", chatResp.Error.Message, chatResp.Error.Type)
	}

	return &chatResp, nil
}

// ChatStream sends a chat request to the LLM and streams the response
func (c *LLMClient) ChatStream(ctx context.Context, messages []ChatMessage, tools []Tool) (<-chan LLMChatStreamResponse, error) {
	req := LLMChatRequest{
		Model:       c.model,
		Messages:    messages,
		Tools:       tools,
		Temperature: 0.7,
		Stream:      true,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("LLM API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	stream := make(chan LLMChatStreamResponse)

	go func() {
		defer resp.Body.Close()
		defer close(stream)

		reader := bufio.NewReader(resp.Body)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						return
					}
					return
				}

				line = strings.TrimSpace(line)
				if !strings.HasPrefix(line, "data:") {
					continue
				}

				data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				if data == "[DONE]" {
					return
				}

				var streamResp LLMChatStreamResponse
				if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
					stream <- LLMChatStreamResponse{Error: &struct {
						Message string `json:"message"`
						Type    string `json:"type"`
						Code    string `json:"code"`
					}{Message: fmt.Sprintf("Failed to unmarshal: %v", err)}}
					return
				}
				stream <- streamResp
			}
		}
	}()

	return stream, nil
}
