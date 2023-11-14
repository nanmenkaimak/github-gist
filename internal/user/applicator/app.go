package applicator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nanmenkaimak/github-gist/internal/user/config"
	"github.com/nanmenkaimak/github-gist/internal/user/controller/grpc"
	"github.com/nanmenkaimak/github-gist/internal/user/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/user/repository"
	"go.uber.org/zap"
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

	grpcService := grpc.NewService(l, repo)
	grpcServer := grpc.NewServer(cfg.GrpcServer.Port, grpcService)

	err = grpcServer.Start()
	if err != nil {
		l.Panicf("failed to start grpc-server err: %v", err)
	}

	defer grpcServer.Close()

	a.gracefulShutdown(cancel)
}

func (a *App) gracefulShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	cancel()
}
