package memory

import (
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
)

type ShortTermMemory struct {
	logger *zap.Logger
	config *pkg.Config
	data   *data.Data
}

func NewShortTermMemory(logger *zap.Logger, config *pkg.Config, data *data.Data) (*ShortTermMemory, error) {
	return &ShortTermMemory{
		logger: logger,
		config: config,
		data:   data,
	}, nil
}

func (s *ShortTermMemory) Search() {

}

func (s *ShortTermMemory) Save() {

}
