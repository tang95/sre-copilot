package memory

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

type Contextual struct {
	ltm     *LongTermMemory
	stm     *ShortTermMemory
	prompts *pkg.PromptConfig
	logger  *zap.Logger
	data    *data.Data
}

func NewContextual(logger *zap.Logger, cfg *pkg.Config, data *data.Data) (*Contextual, error) {
	ltm, err := NewLongTermMemory(logger, cfg, data)
	if err != nil {
		return nil, err
	}
	stm, err := NewShortTermMemory(logger, cfg, data)
	if err != nil {
		return nil, err
	}
	return &Contextual{
		ltm:     ltm,
		stm:     stm,
		prompts: &cfg.Prompt,
		logger:  logger,
		data:    data,
	}, nil
}

func (c *Contextual) BuildMessages(ctx context.Context, userID, chatID, input string) ([]openai.ChatCompletionMessage, error) {
	systemMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: c.prompts.Main,
	}
	messages := []openai.ChatCompletionMessage{systemMessage}
	// 从长期记忆中获取消息
	longTermMessages, err := c.fetchLongTermMemory(ctx, userID, chatID, input)
	if err != nil {
		return nil, err
	}
	if longTermMessages == nil && len(longTermMessages) > 0 {
		messages = append(messages, longTermMessages...)
	}
	// 从短期记忆中获取消息
	shortTermMessages, err := c.fetchShortTermMemory(ctx, userID, chatID, input)
	if err != nil {
		return nil, err
	}
	if shortTermMessages == nil && len(shortTermMessages) > 0 {
		messages = append(messages, shortTermMessages...)
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})
	return messages, nil
}

func (c *Contextual) fetchLongTermMemory(ctx context.Context, userID, chatID, input string) ([]openai.ChatCompletionMessage, error) {
	return nil, nil
}

func (c *Contextual) fetchShortTermMemory(ctx context.Context, userID, chatID, input string) ([]openai.ChatCompletionMessage, error) {
	return nil, nil
}

func (c *Contextual) Save(ctx context.Context, userID, chatID string, messages ...openai.ChatCompletionMessage) error {
	return nil
}
