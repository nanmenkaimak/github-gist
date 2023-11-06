package applicator

import (
	"context"
	"github.com/nanmenkaimak/github-gist/internal/gist/auth"
	"github.com/nanmenkaimak/github-gist/internal/gist/config"
	http2 "github.com/nanmenkaimak/github-gist/internal/gist/controller/http"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"github.com/nanmenkaimak/github-gist/internal/gist/repository"
	"github.com/nanmenkaimak/github-gist/internal/gist/transport"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	logger *zap.SugaredLogger
	config *config.Config
}

func NewApp(logger *zap.SugaredLogger, config *config.Config) *App {
	return &App{
		logger: logger,
		config: config,
	}
}

func (a *App) Run() {
	var cfg = a.config
	var l = a.logger

	ctx, cancel := context.WithCancel(context.TODO())
	_ = ctx

	mainDB, err := dbpostgres.New(cfg.Database.Main)
	if err != nil {
		l.Fatalf("cannot сonnect to mainDB '%s:%d': %v", cfg.Database.Main.Host, cfg.Database.Main.Port, err)
	}

	defer func() {
		if err := mainDB.Close(); err != nil {
			l.Panicf("failed to close mainDB err: %v", err)
		}
	}()
	replicaDB, err := dbpostgres.New(cfg.Database.Replica)
	if err != nil {
		l.Fatalf("cannot сonnect to replicaDB '%s:%d': %v", cfg.Database.Replica.Host, cfg.Database.Replica.Port, err)
	}

	defer func() {
		if err := replicaDB.Close(); err != nil {
			l.Panicf("failed to close replicaDB err: %v", err)
		}
	}()

	repo := repository.NewRepository(mainDB, replicaDB)

	authService := auth.NewService(cfg.Auth)

	userGrpcTransport := transport.NewUserGrpcTransport(cfg.Transport.UserGrpc)

	gistService := gist.NewGistService(repo, userGrpcTransport)

	authMiddleware := middleware.NewJwtV1Middleware(authService, l)

	endpointHandler := http2.NewEndpointHandler(gistService, l)

	router := http2.NewRouter(l, authMiddleware)
	httpCfg := cfg.HttpServer
	server, err := http2.NewServer(httpCfg.Port, httpCfg.ShutdownTimeout, router, l, endpointHandler)
	if err != nil {
		l.Fatalf("failed to create server err: %v", err)
	}

	server.Run()
	defer func() {
		if err := server.Stop(); err != nil {
			l.Panicf("failed close server err: %v", err)
		}
		l.Info("server closed")
	}()

	a.gracefulShutdown(cancel)
}

func (a *App) gracefulShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	cancel()
}
