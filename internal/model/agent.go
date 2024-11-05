package model

import "context"

type AgentMessage struct {
	BaseModel
	ChatID  string `gorm:"type:varchar(36); not null; comment:会话ID" json:"chat_id"`
	UserID  string `gorm:"type:varchar(36); not null; comment:用户ID" json:"user_id"`
	Message string `gorm:"type:text; comment:消息" json:"message"`
}

func (AgentMessage) TableName() string {
	return "agent_message"
}

type AgentMessageRepo interface {
	QueryByChatID(ctx context.Context, chatID string) ([]*AgentMessage, error)
	BatchCreate(ctx context.Context, messages []*AgentMessage) error
	ClearByChatID(ctx context.Context, chatID string) error
}
