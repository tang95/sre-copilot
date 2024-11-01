package memory

import (
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

type LongTermMemory struct {
	logger *zap.Logger
	config *pkg.Config
	data   *data.Data
}

func NewLongTermMemory(logger *zap.Logger, config *pkg.Config, data *data.Data) (*LongTermMemory, error) {
	return &LongTermMemory{
		logger: logger,
		config: config,
		data:   data,
	}, nil
}
