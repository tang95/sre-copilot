package memory

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

type LongTermMemory struct {
	logger *zap.Logger
	config *pkg.Config
	data   *data.Data
	model  *openai.Client
}

func NewLongTermMemory(logger *zap.Logger, config *pkg.Config, data *data.Data) (*LongTermMemory, error) {
	openAiClient, err := pkg.CreateOpenAiClient(&config.Model)
	if err != nil {
		return nil, err
	}
	return &LongTermMemory{
		logger: logger,
		config: config,
		data:   data,
		model:  openAiClient,
	}, nil
}

func (ltm *LongTermMemory) Save(ctx context.Context, userID, chatID string, messages, newMessages []openai.ChatCompletionMessage, totalTokens, newTotalTokens int) error {
	return nil
}
