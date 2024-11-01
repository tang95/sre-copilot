package server

import (
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/tang95/sre-copilot/internal/controller"
	"github.com/tang95/sre-copilot/pkg"
	"github.com/tang95/sre-copilot/pkg/middleware"
	"go.uber.org/zap"
	"time"
)

type HttpServer struct {
	ginEngine *gin.Engine
	logger    *zap.Logger
	config    *pkg.Config
}

func (server *HttpServer) Start() error {
	server.logger.Info(fmt.Sprintf("http server start, addr: %s", server.config.Http.Addr))
	err := server.ginEngine.Run(server.config.Http.Addr)
	if err != nil {
		return err
	}
	return nil
}

func (server *HttpServer) Stop() error {
	server.logger.Info("http server stop")
	return nil
}

func NewHttpServer(config *pkg.Config, logger *zap.Logger, _ *controller.Controller) (pkg.Server, error) {
	if config.Http.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	ginEngine := gin.New()
	ginEngine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	ginEngine.Use(ginzap.RecoveryWithZap(logger, true))
	ginEngine.Use(middleware.Timeout(config))
	return &HttpServer{
		ginEngine: ginEngine,
		logger:    logger,
		config:    config,
	}, nil
}
