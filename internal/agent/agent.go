package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/tang95/sre-copilot/internal/agent/memory"
	"github.com/tang95/sre-copilot/internal/agent/tool"
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/internal/service"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

type Agent struct {
	model  *openai.Client
	config *pkg.Config
	logger *zap.Logger
	svc    *service.Service
	data   *data.Data
	memory *memory.Contextual
}

func NewAgent(logger *zap.Logger, cfg *pkg.Config, svc *service.Service, d *data.Data) (*Agent, error) {
	openAiClient, err := pkg.CreateOpenAiClient(&cfg.Model)
	if err != nil {
		return nil, err
	}
	contextual, err := memory.NewContextual(logger, cfg, d)
	if err != nil {
		return nil, err
	}
	return &Agent{
		model:  openAiClient,
		config: cfg,
		logger: logger,
		svc:    svc,
		data:   d,
		memory: contextual,
	}, nil
}

func (a *Agent) chatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage, tools []openai.Tool) (*openai.ChatCompletionResponse, *openai.ChatCompletionMessage, error) {
	response, err := a.model.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    a.config.Model.Model,
		Messages: messages,
		Tools:    tools,
	})
	if err != nil {
		return nil, nil, err
	}
	if len(response.Choices) != 1 {
		return nil, nil, errors.New("expected 1 choice")
	}
	return &response, &response.Choices[0].Message, nil
}

func (a *Agent) buildTools(_ context.Context) ([]openai.Tool, error) {
	tools := make([]openai.Tool, 0)
	for i := range tool.Tools {
		t := tool.Tools[i]
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  t.Parameters(),
			},
		})
	}
	return tools, nil
}

func (a *Agent) handleToolCalls(ctx context.Context, toolCalls []openai.ToolCall) ([]openai.ChatCompletionMessage, error) {
	toolCallMessages := make([]openai.ChatCompletionMessage, 0)
	tools := make(map[string]tool.Tool)
	for _, t := range tool.Tools {
		tools[t.Name()] = t
	}

	for _, toolCall := range toolCalls {
		t, ok := tools[toolCall.Function.Name]
		if !ok {
			toolCallMessages = append(toolCallMessages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    "tool not found",
				Name:       toolCall.Function.Name,
				ToolCallID: toolCall.ID,
			})
			continue
		}

		toolResult, err := t.Call(ctx, toolCall.Function.Arguments)
		if err != nil {
			toolCallMessages = append(toolCallMessages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    fmt.Sprintf("failed to call tool result: %s error: %v", toolResult, err),
				Name:       toolCall.Function.Name,
				ToolCallID: toolCall.ID,
			})
		} else {
			toolCallMessages = append(toolCallMessages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    toolResult,
				Name:       toolCall.Function.Name,
				ToolCallID: toolCall.ID,
			})
		}
	}
	return toolCallMessages, nil
}

func (a *Agent) Invoke(ctx context.Context, userID, chatID, input string) ([]openai.ChatCompletionMessage, error) {
	tools, err := a.buildTools(ctx)
	if err != nil {
		return nil, err
	}
	newMessages := make([]openai.ChatCompletionMessage, 0)
	newMessages = append(newMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})
	messages, err := a.memory.BuildMessages(ctx, userID, chatID, input)
	if err != nil {
		return nil, err
	}
	resp, msg, err := a.chatCompletion(ctx, append(messages, newMessages...), tools)
	if err != nil {
		return nil, err
	}
	newMessages = append(newMessages, *msg)
	messageTokens := resp.Usage.PromptTokens
	for {
		if len(msg.ToolCalls) > 0 {
			toolCallMessages, err := a.handleToolCalls(ctx, msg.ToolCalls)
			if err != nil {
				return nil, err
			}
			newMessages = append(newMessages, toolCallMessages...)
			resp, msg, err = a.chatCompletion(ctx, append(messages, newMessages...), tools)
			if err != nil {
				return nil, err
			}
			newMessages = append(newMessages, *msg)
		} else {
			break
		}
	}
	err = a.memory.Save(ctx, userID, chatID, messages, newMessages, resp.Usage.TotalTokens, resp.Usage.TotalTokens-messageTokens)
	if err != nil {
		return nil, err
	}
	return newMessages, nil
}
