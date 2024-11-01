package data

import (
	"context"
	"errors"
	"github.com/tang95/sre-copilot/internal/model"
)

type agentStateRepo struct {
	*Data
}

func newAgentStateRepo(data *Data) model.AgentStateRepo {
	return &agentStateRepo{data}
}

func (repo *agentStateRepo) Get(ctx context.Context, chatID string) (*model.AgentState, error) {
	state := model.AgentState{}
	tx := repo.DB(ctx).Where("chat_id = ?", chatID).First(&state)
	return &state, tx.Error
}

func (repo *agentStateRepo) Create(ctx context.Context, state *model.AgentState) error {
	tx := repo.DB(ctx).Create(&state)
	return tx.Error
}

func (repo *agentStateRepo) Update(ctx context.Context, state *model.AgentState) error {
	if state.ChatID == "" {
		return errors.New("state must have a chat ID")
	}
	tx := repo.DB(ctx).Where("chat_id =?", state.ChatID).Updates(&state)
	return tx.Error
}
