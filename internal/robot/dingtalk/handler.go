package dingtalk

import (
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"go.uber.org/zap"
)

func (d *DingTalk) OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	replier := chatbot.NewChatbotReplier()
	messages, err := d.agent.Invoke(ctx, data.SenderId, data.ConversationId, data.Text.Content)
	if err != nil {
		d.logger.Error("invoke agent error", zap.Error(err))
		return nil, err
	}
	if err := replier.SimpleReplyText(ctx, data.SessionWebhook, []byte(messages[len(messages)-1].Content)); err != nil {
		d.logger.Error("reply error", zap.Error(err))
		return nil, err
	}
	return []byte("ok"), nil
}
