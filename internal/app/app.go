package app

import (
	"context"
	"doc-server/internal/api"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type App struct {
	serviceProvider *service_provider
	router          *gin.Engine
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	err := a.runHttpServer()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runHttpServer() error {
	zap.L().Info(fmt.Sprintf("HTTP server is running on %s", a.serviceProvider.ConfigHTTP().Address()))

	err := a.router.Run(a.serviceProvider.ConfigHTTP().Address())
	return err
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHttpServer,
		a.initGlobalLogger,
	}

	for _, init := range inits {
		err := init(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = NewServiceProvider()
	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api.RegisterRoutes(router, a.serviceProvider.AuthHandler(ctx), a.serviceProvider.DocsHandler(ctx))

	a.router = router
	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "local.env", "config path")
	flag.Parse()
	godotenv.Load(configPath)

	return nil
}

func (a *App) initGlobalLogger(ctx context.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	// Устанавливаем логгер глобальным
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
	return nil
}
