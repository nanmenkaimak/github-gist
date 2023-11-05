package applicator

import (
	"context"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
	"github.com/nanmenkaimak/github-gist/internal/auth/config"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/consumer"
	http2 "github.com/nanmenkaimak/github-gist/internal/auth/controller/http"
	"github.com/nanmenkaimak/github-gist/internal/auth/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/auth/database/dbredis"
	"github.com/nanmenkaimak/github-gist/internal/auth/repository"
	"github.com/nanmenkaimak/github-gist/internal/auth/transport"
	"github.com/nanmenkaimak/github-gist/internal/kafka"
	"go.uber.org/zap"
	"net/http"
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

	dbRedis, err := dbredis.New()
	if err != nil {
		l.Fatalf("cannot connect to redis")
	}

	userVerificationProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		l.Panicf("failed NewProducer err: %v", err)
	}

	userVerificationConsumerCallback := consumer.NewUserVerificationCallback(l)

	userVerificationConsumer, err := kafka.NewConsumer(l, cfg.Kafka, userVerificationConsumerCallback)
	if err != nil {
		l.Panicf("failed NewConsumer err: %v", err)
	}

	go userVerificationConsumer.Start()

	repo := repository.NewRepository(mainDB, replicaDB)

	client := &http.Client{}

	userTransport := transport.NewTransport(cfg.Transport.User, client)

	authService := auth.NewAuthService(repo, cfg.Auth, userTransport, userVerificationProducer, dbRedis)

	endpointHandler := http2.NewEndpointHandler(authService, l)

	router := http2.NewRouter(l)
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
