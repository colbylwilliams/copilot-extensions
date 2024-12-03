package copilot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/colbylwilliams/copilot-extensions/pkg/jsonschema"
)

type Client struct {
	// httpClient is the HTTP client used to communicate with the API.
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}

type ChatMessage struct {
	Role          string                   `json:"role"`
	Content       string                   `json:"content"`
	Name          string                   `json:"name,omitempty"`
	FunctionCall  *ChatMessageFunctionCall `json:"function_call,omitempty"`
	Functions     []*ChatMessageFunction   `json:"functions,omitempty"`
	Confirmations []*ChatConfirmation      `json:"copilot_confirmations"`
	ToolCalls     []*ToolCall              `json:"tool_calls"`
}

type ToolCall struct {
	Function *ChatMessageFunctionCall
}

type ChatMessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatMessageFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  *jsonschema.Definition `json:"parameters"`
}

type ChatConfirmation struct {
	State        string `json:"state"`
	Confirmation any    `json:"confirmation"`
}

type ChatCompletionResponse struct {
	Choices []ChatChoice `json:"choices"`
}

type ChatChoice struct {
	Message ChatMessage `json:"message"`
}

type ChatModel string

const (
	ChatModelGPT35 ChatModel = "gpt-3.5-turbo"
	ChatModelGPT4  ChatModel = "gpt-4"
)

type ChatCompletionsRequest struct {
	Messages []ChatMessage  `json:"messages"`
	Model    ChatModel      `json:"model"`
	Stream   bool           `json:"stream"`
	Tools    []FunctionTool `json:"tools"`
}

type FunctionTool struct {
	Type     string              `json:"type"`
	Function ChatMessageFunction `json:"function"`
}

func (c *Client) ChatCompletionsStream(ctx context.Context, integrationID, token string, req ChatCompletionsRequest) (io.ReadCloser, error) {
	req.Stream = true
	return c.ChatCompletions(ctx, integrationID, token, req)
}

func (c *Client) ChatCompletions(ctx context.Context, integrationID, token string, req ChatCompletionsRequest) (io.ReadCloser, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.githubcopilot.com/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	if integrationID != "" {
		httpReq.Header.Set("Copilot-Integration-Id", integrationID)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		fmt.Println(string(b))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
