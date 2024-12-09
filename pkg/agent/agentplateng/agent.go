package agentplateng

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/colbylwilliams/copilot-extensions/pkg/agent"
	"github.com/colbylwilliams/copilot-extensions/pkg/github"
	oai "github.com/colbylwilliams/copilot-extensions/pkg/openai"
	"github.com/colbylwilliams/copilot-go"
	"github.com/colbylwilliams/copilot-go/sse"
	"github.com/openai/openai-go"
)

type Agent struct {
	cfg *copilot.Config
	oai *openai.Client
}

func New(cfg *copilot.Config, oai *openai.Client) *Agent {
	return &Agent{
		cfg: cfg,
		oai: oai,
	}
}

func (a *Agent) Route() string {
	return "plat-eng"
}

func (a *Agent) Execute(ctx context.Context, token string, req *copilot.Request, w http.ResponseWriter) error {

	// get the github user
	gh := github.GetUserClient(a.cfg, token)
	me, _, err := gh.Users.Get(ctx, "")
	if err != nil {
		return err
	}
	fmt.Println("me: ", me.GetName())

	// write the sse headers
	sse.WriteStreamingHeaders(w)

	messages := []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(PromptStart)}

	session := copilot.GetSessionInfo(ctx)
	if session != nil {
		agent.PrintJsonWithNewLines(session)
	}

	for _, m := range req.Messages {
		// for now, we won't send the _session message
		// web (dotcom) chat sends us downstream to openai
		if m.IsSessionMessage() {
			continue
		}

		// we'll also skip adding messages with no content
		if m.Content == "" {
			continue
		}

		switch m.Role {
		case copilot.ChatRoleSystem:
			messages = append(messages, openai.SystemMessage(m.Content))

		case copilot.ChatRoleUser:
			// if the message begins with @agent-name then remove it
			messages = append(messages, openai.UserMessage(strings.TrimPrefix(m.Content, fmt.Sprintf("@%s ", req.Agent))))

		case copilot.ChatRoleAssistant:
			messages = append(messages, openai.AssistantMessage(m.Content))

		default:
			return fmt.Errorf("unhandled role: %s", m.Role)
		}
	}

	// write the request to a file
	// if err := agent.WriteRequestToFile(ctx, *req); err != nil {
	// 	return err
	// }

	// write the response to a file
	// f, err := agent.GetTempLogFileForRequest(ctx)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()

	stream := a.oai.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(a.cfg.ChatModel),
	})

	for stream.Next() {
		evt := stream.Current()
		chunk, err := oai.ToAgentResponse(&evt)
		if err != nil {
			return err
		}

		sse.WriteData(w, chunk)
		// sse.WriteData(f, chunk)
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return nil
}
