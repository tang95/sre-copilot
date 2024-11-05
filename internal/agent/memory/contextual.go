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

func (contextual *Contextual) BuildMessages(ctx context.Context, userID, chatID, input string) ([]openai.ChatCompletionMessage, error) {
	messages := make([]openai.ChatCompletionMessage, 0)
	// 从长期记忆中获取消息
	longTermMessages, err := contextual.fetchLongTermMemory(ctx, userID, chatID, input)
	if err != nil {
		return nil, err
	}
	if longTermMessages != nil && len(longTermMessages) > 0 {
		messages = append(messages, longTermMessages...)
	}
	// 从短期记忆中获取消息
	shortTermMessages, err := contextual.fetchShortTermMemory(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}
	if shortTermMessages != nil && len(shortTermMessages) > 0 {
		messages = append(messages, shortTermMessages...)
	}
	return messages, nil
}

func (contextual *Contextual) fetchLongTermMemory(ctx context.Context, userID, chatID, input string) ([]openai.ChatCompletionMessage, error) {
	// TODO: 从长期记忆获取信息拼接SystemMessage
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: contextual.prompts.Main,
		},
	}, nil
}

func (contextual *Contextual) fetchShortTermMemory(ctx context.Context, userID, chatID string) ([]openai.ChatCompletionMessage, error) {
	messages, err := contextual.stm.Search(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (contextual *Contextual) Save(ctx context.Context, userID, chatID string, messages, newMessages []openai.ChatCompletionMessage, totalTokens, newTotalTokens int) error {
	err := contextual.stm.Save(ctx, userID, chatID, messages, newMessages, totalTokens, newTotalTokens)
	if err != nil {
		return err
	}
	err = contextual.ltm.Save(ctx, userID, chatID, messages, newMessages, totalTokens, newTotalTokens)
	if err != nil {
		return err
	}
	return nil
}
