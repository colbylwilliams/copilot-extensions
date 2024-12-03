package agentplateng

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/colbylwilliams/copilot-extensions/pkg/agent"
	"github.com/colbylwilliams/copilot-extensions/pkg/config"
	"github.com/colbylwilliams/copilot-extensions/pkg/github"

	"github.com/colbylwilliams/copilot-extensions/pkg/convert"
	"github.com/colbylwilliams/copilot-extensions/pkg/sse"
	"github.com/openai/openai-go"
	// octokit "github.com/octokit/go-sdk/pkg"
)

type Agent struct {
	cfg *config.Config
	oai *openai.Client
}

func New(cfg *config.Config, oai *openai.Client) *Agent {
	return &Agent{
		cfg: cfg,
		oai: oai,
	}
}

func (a *Agent) Execute(ctx context.Context, integrationID, apiToken string, req *agent.Request, w http.ResponseWriter) error {

	// get the github user
	gh := github.GetUserClient(a.cfg, apiToken)
	me, _, err := gh.Users.Get(ctx, "")
	if err != nil {
		return err
	}
	fmt.Println("me: ", me.GetName())

	// flusher, ok := w.(http.Flusher)
	// if !ok {
	// 	return fmt.Errorf("sse not supported")
	// }

	// write the sse headers
	sse.WriteStreamingHeaders(w)

	messages := []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(PromptStart)}

	session, err := req.GetSessionContext()
	if err != nil {
		fmt.Println("error: ", err)
	}

	if session != nil {
		agent.PrintJsonWithNewLines(session)

		if session.Item != nil {
			switch session.Item.Type {
			case agent.RepoItemRefTypeIssue:
				issue, _, err := gh.Issues.Get(ctx, session.Item.Owner, session.Item.Repo, session.Item.Number)
				if err != nil {
					return fmt.Errorf("failed to get issue: %w", err)
				}
				fmt.Println("issue: ", issue)

				vars, _, err := gh.Actions.ListRepoVariables(ctx, session.Item.Owner, session.Item.Repo, nil)
				if err != nil {
					return fmt.Errorf("failed to get repo variables: %w", err)
				}

				fmt.Println("vars: ", vars)
				for _, v := range vars.Variables {
					fmt.Printf("  %s: %s\n", v.Name, v.Value)
				}

			case agent.RepoItemRefTypePull:
				pr, _, err := gh.PullRequests.Get(ctx, session.Item.Owner, session.Item.Repo, session.Item.Number)
				if err != nil {
					return fmt.Errorf("failed to get pull request: %w", err)
				}
				fmt.Println("pr: ", pr)

				vars, _, err := gh.Actions.ListRepoVariables(ctx, session.Item.Owner, session.Item.Repo, nil)
				if err != nil {
					return fmt.Errorf("failed to get repo variables: %w", err)
				}

				fmt.Println("vars: ", vars)
				for _, v := range vars.Variables {
					fmt.Printf("  %s: %s\n", v.Name, v.Value)
				}

			default:
				fmt.Println("warning: unhandled session item type")
			}
		}
	}

	for _, m := range req.Messages {
		// for now, we won't send the _session message
		// web (dotcom) chat sends us downstream to openai
		if m.IsSession() {
			continue
		}

		// we'll also skip adding messages with no content
		if m.Content == "" {
			continue
		}

		switch m.Role {
		case agent.ChatRoleSystem:
			messages = append(messages, openai.SystemMessage(m.Content))

		case agent.ChatRoleUser:
			messages = append(messages, openai.UserMessage(
				// if the message begins with @agent-name then remove it
				strings.TrimPrefix(m.Content, fmt.Sprintf("@%s ", req.Agent))))

		case agent.ChatRoleAssistant:
			messages = append(messages, openai.AssistantMessage(m.Content))

		default:
			return fmt.Errorf("invalid role: %s", m.Role)
		}
	}

	// write the request to a file
	if err := agent.WriteRequestToFile(ctx, *req); err != nil {
		return err
	}

	// write the response to a file
	f, err := agent.GetTempLogFileForRequest(ctx)
	if err != nil {
		return err
	}
	defer f.Close()

	stream := a.oai.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(a.cfg.ChatModel),
	})

	for stream.Next() {
		evt := stream.Current()
		chunk, err := convert.ToAgentResponse(&evt)
		if err != nil {
			return err
		}

		sse.WriteDataAndFlush(w, chunk)
		sse.WriteData(f, chunk)
		// flusher.Flush()
	}

	// sse.WriteDone(w)
	// flusher.Flush()
	// sse.WriteDone(f)

	if err := stream.Err(); err != nil {
		return err
	}

	return nil
}
