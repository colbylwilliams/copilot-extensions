package openai

import (
	"github.com/colbylwilliams/copilot-go"
	"github.com/openai/openai-go"
)

// agentResponse converts a openai ChatCompletionChunk
// to the copilot representation of agent.Response.
func ToAgentResponse(in *openai.ChatCompletionChunk) (*copilot.Response, error) {
	out := copilot.Response{
		ID:                in.ID,
		Created:           in.Created,
		Object:            string(in.Object),
		Model:             in.Model,
		SystemFingerprint: in.SystemFingerprint,
		Choices:           make([]copilot.ChatChoice, len(in.Choices)),
	}

	for i, c := range in.Choices {
		newChoice := copilot.ChatChoice{
			Index:        c.Index,
			FinishReason: string(c.FinishReason),
			Delta: copilot.ChatChoiceDelta{
				Role:    string(c.Delta.Role),
				Content: c.Delta.Content,
			},
		}
		out.Choices[i] = newChoice
	}

	return &out, nil
}
