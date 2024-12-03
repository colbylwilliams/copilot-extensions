package convert

import (
	"github.com/colbylwilliams/copilot-extensions/pkg/agent"
	"github.com/openai/openai-go"
)

// https://github.com/github/copilot-api/pkg/azure/azuremodel/models.go

// ToAgentChatMessage converts a ToAgentChatMessage from the openai sdk to the
// internal representation of openai.ChatCompletions.
func ToAgentResponse(in *openai.ChatCompletionChunk) (*agent.Response, error) {
	out := agent.Response{
		ID:                in.ID,
		Created:           in.Created,
		Object:            string(in.Object),
		Model:             in.Model,
		SystemFingerprint: in.SystemFingerprint,
		Choices:           make([]agent.ChatChoice, len(in.Choices)),
		// Usage:               FromAzureCompletionsUsage(in.Usage),
		// PromptFilterResults: FromAzurePromptFilterResults(in.PromptFilterResults),
	}

	for i, c := range in.Choices {
		newChoice := agent.ChatChoice{
			Index:        c.Index,
			FinishReason: string(c.FinishReason),
			Delta:        ToAgentChatChoiceDelta(c.Delta),
		}
		out.Choices[i] = newChoice
	}

	return &out, nil
}

// ToAgentChatChoiceDelta converts a ChatChoiceDelta from the openai sdk to the
// internal representation of agent.ChatChoiceDelta.
func ToAgentChatChoiceDelta(delta openai.ChatCompletionChunkChoicesDelta) agent.ChatChoiceDelta {
	return agent.ChatChoiceDelta{
		Role:    string(delta.Role),
		Content: delta.Content,
		// seems copilot chokes on this
		// FunctionCall: ChatCompletionChunkChoicesDeltaFunctionCall(&delta.FunctionCall),
		// Context:      FromAzureChatMessageContext(delta.Context),
		// ToolCalls:    FromAzureChatMessageToolCalls(delta.ToolCalls),
	}
}

// ToAgentChoiceDeltaFunctionCall converts a ChatMessageFunctionCall from the openai sdk to the
// internal representation of agent.ChatMessageFunctionCall.
func ToAgentChoiceDeltaFunctionCall(functionCall *openai.ChatCompletionChunkChoicesDeltaFunctionCall) *agent.ChatChoiceDeltaFunctionCall {
	if functionCall == nil {
		return nil
	}
	return &agent.ChatChoiceDeltaFunctionCall{
		Name:      functionCall.Name,
		Arguments: functionCall.Arguments,
	}
}
