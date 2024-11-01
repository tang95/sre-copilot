package model

import (
	"context"
)

type AgentState struct {
	BaseModel
	ChatID   string `gorm:"type:varchar(36); not null; comment:会话ID" json:"chat_id"`
	Messages string `gorm:"type:text; comment:消息" json:"messages"`
}

func (AgentState) TableName() string {
	return "agent_state"
}

type AgentStateRepo interface {
	Get(ctx context.Context, chatID string) (*AgentState, error)
	Create(ctx context.Context, state *AgentState) error
	Update(ctx context.Context, state *AgentState) error
}
