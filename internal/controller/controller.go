package controller

import (
	"embed"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/tang95/sre-copilot/internal/service"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

//go:embed static
var consoleFS embed.FS

type Controller struct {
	service     *service.Service
	config      *pkg.Config
	logger      *zap.Logger
	transaction service.Transaction
}

func NewController(
	service *service.Service,
	config *pkg.Config,
	logger *zap.Logger,
	transaction service.Transaction,
) (*Controller, error) {
	return &Controller{
		service:     service,
		config:      config,
		logger:      logger,
		transaction: transaction,
	}, nil
}

func (controller *Controller) WithRoutes(engine *gin.Engine, jwtMiddleware *jwt.GinJWTMiddleware) {
	// console
	consoleServer := static.Serve("/", static.EmbedFolder(consoleFS, "static"))
	engine.Use(consoleServer)
	engine.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodGet &&
			!strings.ContainsRune(ctx.Request.URL.Path, '.') &&
			!strings.HasPrefix(ctx.Request.URL.Path, "/api/") {
			ctx.Request.URL.Path = "/"
			consoleServer(ctx)
		}
	})
	// api group
	api := engine.Group("/api")

	// component
	component := api.Group("/incident", jwtMiddleware.MiddlewareFunc())
	component.GET("/query", controller.queryIncidents())
}
