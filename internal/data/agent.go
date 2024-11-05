package data

import (
	"context"
	"github.com/tang95/sre-copilot/internal/model"
)

type agentMessageRepo struct {
	*Data
}

func newAgentMessageRepo(data *Data) model.AgentMessageRepo {
	return &agentMessageRepo{data}
}

func (repo *agentMessageRepo) QueryByChatID(ctx context.Context, chatID string) ([]*model.AgentMessage, error) {
	messages := make([]*model.AgentMessage, 0)
	tx := repo.DB(ctx).Where("chat_id =?", chatID).Find(&messages)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return messages, nil
}

func (repo *agentMessageRepo) BatchCreate(ctx context.Context, messages []*model.AgentMessage) error {
	tx := repo.DB(ctx).Create(&messages)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (repo *agentMessageRepo) ClearByChatID(ctx context.Context, chatID string) error {
	tx := repo.DB(ctx).Where("chat_id =?", chatID).Delete(&model.AgentMessage{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
