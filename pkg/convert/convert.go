package convert

import (
	"github.com/colbylwilliams/copilot-extensions/pkg/agent"
	"github.com/openai/openai-go"
)

// ToAgentResponse converts a openai ChatCompletionChunk
// to the copilot representation of agent.Response.
func ToAgentResponse(in *openai.ChatCompletionChunk) (*agent.Response, error) {
	out := agent.Response{
		ID:                in.ID,
		Created:           in.Created,
		Object:            string(in.Object),
		Model:             in.Model,
		SystemFingerprint: in.SystemFingerprint,
		Choices:           make([]agent.ChatChoice, len(in.Choices)),
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

// ToAgentChatChoiceDelta converts a openai ChatCompletionChunkChoicesDelta
// to the copilot representation of agent.ChatChoiceDelta.
func ToAgentChatChoiceDelta(delta openai.ChatCompletionChunkChoicesDelta) agent.ChatChoiceDelta {
	return agent.ChatChoiceDelta{
		Role:    string(delta.Role),
		Content: delta.Content,
	}
}

// ToAgentChoiceDeltaFunctionCall converts a openai ChatCompletionChunkChoicesDeltaFunctionCall
// to the copilot representation of agent.ChatChoiceDeltaFunctionCall
func ToAgentChoiceDeltaFunctionCall(functionCall *openai.ChatCompletionChunkChoicesDeltaFunctionCall) *agent.ChatChoiceDeltaFunctionCall {
	if functionCall == nil {
		return nil
	}
	return &agent.ChatChoiceDeltaFunctionCall{
		Name:      functionCall.Name,
		Arguments: functionCall.Arguments,
	}
}
