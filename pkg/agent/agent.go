package agent

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleFunction  ChatRole = "function"
	ChatRoleTool      ChatRole = "tool"
	ChatRoleSystem    ChatRole = "system"
)

// Request is the request to the agent.
type Request struct {
	ThreadID         string     `json:"copilot_thread_id"`
	Messages         []*Message `json:"messages"`
	Stop             []string   `json:"stop"`
	TopP             float32    `json:"top_p"`
	Temperature      float32    `json:"temperature"`
	MaxTokens        int32      `json:"max_tokens"`
	PresencePenalty  float32    `json:"presence_penalty"`
	FrequencyPenalty float32    `json:"frequency_penalty"`
	Skills           []string   `json:"copilot_skills"`
	Agent            string     `json:"agent"`
}

// Message is a message in the request.
type Message struct {
	Role          ChatRole              `json:"role"`
	Content       string                `json:"content"`
	Name          string                `json:"name,omitempty"`
	References    []*Reference          `json:"copilot_references"`
	Confirmations []*ClientConfirmation `json:"copilot_confirmations"`
	FunctionCall  *FunctionCall         `json:"functionCall,omitempty"`
	ToolCalls     []*ToolCall           `json:"toolCalls,omitempty"`
	ToolCallID    string                `json:"toolCallID,omitempty"`
}

// Response is the response from the agent.
type Response struct {
	ID                string        `json:"id"`
	Created           int64         `json:"created"`
	Object            string        `json:"object" default:"chat.completion.chunk"`
	Model             string        `json:"model"`
	SystemFingerprint string        `json:"system_fingerprint,omitempty"`
	Choices           []ChatChoice  `json:"choices"`
	References        []*Reference  `json:"copilot_references,omitempty"`
	Confirmation      *Confirmation `json:"copilot_confirmation,omitempty"`
	Errors            []*Error      `json:"copilot_errors,omitempty"`
}

type ChatChoice struct {
	Index        int64           `json:"index"`
	FinishReason string          `json:"finish_reason,omitempty"`
	Delta        ChatChoiceDelta `json:"delta"`
}

type ChatChoiceDelta struct {
	Content      string                       `json:"content"`
	Role         string                       `json:"role,omitempty"`
	Name         string                       `json:"name,omitempty"`
	FunctionCall *ChatChoiceDeltaFunctionCall `json:"function_call,omitempty"`
}

type ChatChoiceDeltaFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID       string
	Type     string
	Function *FunctionCall
	Index    int
}
