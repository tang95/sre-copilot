package dingtalk

import (
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/tang95/sre-copilot/internal/agent"
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/internal/service"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

type DingTalk struct {
	config  *pkg.Config
	logger  *zap.Logger
	service *service.Service
	client  *client.StreamClient
	agent   *agent.Agent
	data    *data.Data
}

func NewDingTalk(logger *zap.Logger, cfg *pkg.Config, svc *service.Service, agents *agent.Agent, d *data.Data) (pkg.Server, error) {
	return &DingTalk{
		config:  cfg,
		logger:  logger,
		service: svc,
		agent:   agents,
		data:    d,
	}, nil
}

func (d *DingTalk) Start() error {
	cli := client.NewStreamClient(client.WithAppCredential(client.NewAppCredentialConfig(
		d.config.Robot.ClientId, d.config.Robot.ClientSecret,
	)))
	cli.RegisterChatBotCallbackRouter(d.OnChatBotMessageReceived)
	if err := cli.Start(context.Background()); err != nil {
		return err
	}
	d.client = cli
	d.logger.Info("start dingtalk stream client", zap.String("client_id", d.config.Robot.ClientId))
	return nil
}

func (d *DingTalk) Stop() error {
	d.client.Close()
	d.logger.Info("stop dingtalk stream client")
	return nil
}
