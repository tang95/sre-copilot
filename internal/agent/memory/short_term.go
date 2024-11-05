package memory

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sashabaranov/go-openai"
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/internal/model"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

const MaxTokens = 1024

type ShortTermMemory struct {
	logger *zap.Logger
	config *pkg.Config
	data   *data.Data
	model  *openai.Client
}

func NewShortTermMemory(logger *zap.Logger, config *pkg.Config, data *data.Data) (*ShortTermMemory, error) {
	openAiClient, err := pkg.CreateOpenAiClient(&config.Model)
	if err != nil {
		return nil, err
	}
	return &ShortTermMemory{
		logger: logger,
		config: config,
		data:   data,
		model:  openAiClient,
	}, nil
}

func (stm *ShortTermMemory) Search(ctx context.Context, _, chatID string) ([]openai.ChatCompletionMessage, error) {
	agentMessages, err := stm.data.AgentMessageRepo.QueryByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	messages := make([]openai.ChatCompletionMessage, 0, len(agentMessages))
	for _, msg := range agentMessages {
		var message openai.ChatCompletionMessage
		err = json.Unmarshal([]byte(msg.Message), &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (stm *ShortTermMemory) Save(ctx context.Context, userID, chatID string, messages, newMessages []openai.ChatCompletionMessage, totalTokens, newTotalTokens int) error {
	insertMessages := make([]openai.ChatCompletionMessage, 0)
	if (totalTokens - newTotalTokens) > MaxTokens {
		historyMessages, err := stm.summaryHistoryMessages(ctx, messages)
		if err != nil {
			return err
		}
		insertMessages = append(insertMessages, historyMessages...)
		err = stm.data.AgentMessageRepo.ClearByChatID(ctx, chatID)
		if err != nil {
			return err
		}
	}
	insertMessages = append(insertMessages, newMessages...)
	// 保存消息
	agentMessages := make([]*model.AgentMessage, 0)
	for _, message := range insertMessages {
		marshal, err := json.Marshal(message)
		if err != nil {
			return err
		}
		agentMessages = append(agentMessages, &model.AgentMessage{
			UserID:  userID,
			ChatID:  chatID,
			Message: string(marshal),
		})
	}
	err := stm.data.AgentMessageRepo.BatchCreate(ctx, agentMessages)
	if err != nil {
		return err
	}
	return nil
}

func (stm *ShortTermMemory) summaryHistoryMessages(ctx context.Context, messages []openai.ChatCompletionMessage) ([]openai.ChatCompletionMessage, error) {
	messages = messages[1:]
	systemPrompt := `Progressively summarize the lines of conversation provided, adding onto the previous summary returning a new summary.

EXAMPLE
Current summary:
The human asks what the AI thinks of artificial intelligence. The AI thinks artificial intelligence is a force for good.

New lines of conversation:
Human: Why do you think artificial intelligence is a force for good?
AI: Because artificial intelligence will help humans reach their full potential.

New summary:
The human asks what the AI thinks of artificial intelligence. The AI thinks artificial intelligence is a force for good because it will help humans reach their full potential.
END OF EXAMPLE
`
	if messages[0].Role == openai.ChatMessageRoleAssistant {
		systemPrompt += messages[0].Content + "\n"
		messages = messages[1:]
	} else {
		systemPrompt += "Current summary：\n\n"
	}
	systemPrompt += "New lines of conversation:\n"
	for _, message := range messages {
		systemPrompt += message.Role + ": " + message.Content + "\n"
		if len(message.ToolCalls) > 0 {
			for _, toolCall := range message.ToolCalls {
				systemPrompt += toolCall.Function.Name + ": " + toolCall.Function.Arguments + "\n"
			}
		}
	}
	systemPrompt += "\nNew summary(use conversation language):\n"
	response, err := stm.model.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: stm.config.Model.Model,
		Messages: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		}},
	})
	if err != nil {
		return nil, err
	}

	if len(response.Choices) != 1 {
		return nil, errors.New("expected 1 choice")
	}
	summary := response.Choices[0].Message
	summary.Content = "Current summary：\n" + summary.Content
	return []openai.ChatCompletionMessage{summary}, nil
}
