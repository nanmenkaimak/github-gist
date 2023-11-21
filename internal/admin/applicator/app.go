package applicator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nanmenkaimak/github-gist/internal/admin/storage"

	"github.com/nanmenkaimak/github-gist/internal/admin/admin"
	"github.com/nanmenkaimak/github-gist/internal/admin/auth"
	"github.com/nanmenkaimak/github-gist/internal/admin/config"
	"github.com/nanmenkaimak/github-gist/internal/admin/controller/http"
	"github.com/nanmenkaimak/github-gist/internal/admin/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/admin/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/admin/repository"
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

	mainDBgist, err := dbpostgres.New(cfg.Databases.Gist.Main)
	if err != nil {
		l.Fatalf("cannot сonnect to mainDB '%s:%d': %v", cfg.Databases.Gist.Main.Host, cfg.Databases.Gist.Main.Port, err)
	}

	defer func() {
		if err := mainDBgist.Close(); err != nil {
			l.Panicf("failed to close mainDB err: %v", err)
		}
	}()
	replicaDBgist, err := dbpostgres.New(cfg.Databases.Gist.Replica)
	if err != nil {
		l.Fatalf("cannot сonnect to replicaDB '%s:%d': %v", cfg.Databases.Gist.Replica.Host, cfg.Databases.Gist.Replica.Port, err)
	}

	defer func() {
		if err := replicaDBgist.Close(); err != nil {
			l.Panicf("failed to close replicaDB err: %v", err)
		}
	}()

	mainDBuser, err := dbpostgres.New(cfg.Databases.User.Main)
	if err != nil {
		l.Fatalf("cannot сonnect to mainDB '%s:%d': %v", cfg.Databases.Gist.Main.Host, cfg.Databases.Gist.Main.Port, err)
	}

	defer func() {
		if err := mainDBuser.Close(); err != nil {
			l.Panicf("failed to close mainDB err: %v", err)
		}
	}()
	replicaDBuser, err := dbpostgres.New(cfg.Databases.User.Replica)
	if err != nil {
		l.Fatalf("cannot сonnect to replicaDB '%s:%d': %v", cfg.Databases.Gist.Replica.Host, cfg.Databases.Gist.Replica.Port, err)
	}

	defer func() {
		if err := replicaDBuser.Close(); err != nil {
			l.Panicf("failed to close replicaDB err: %v", err)
		}
	}()

	userRepo := repository.NewUserRepository(mainDBuser, replicaDBuser)

	gistRepo := repository.NewGistRepository(mainDBgist, replicaDBgist)

	authService := auth.NewService(cfg.Auth)

	userStorage := storage.NewDataStorage(cfg.Storage.Interval, userRepo, l)

	go userStorage.Run()

	adminService := admin.NewAdminService(gistRepo, userRepo, userStorage)

	authMiddleware := middleware.NewJwtV1Middleware(authService, l)

	endpointHandler := http.NewEndpointHandler(adminService, l)

	router := http.NewRouter(l, authMiddleware)

	httpCfg := cfg.HttpServer

	server, err := http.NewServer(httpCfg.Port, httpCfg.ShutdownTimeout, router, l, endpointHandler)
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
