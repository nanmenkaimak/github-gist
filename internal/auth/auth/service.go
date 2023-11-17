package auth

import (
	"github.com/nanmenkaimak/github-gist/internal/auth/config"
	"github.com/nanmenkaimak/github-gist/internal/auth/repository"
	"github.com/nanmenkaimak/github-gist/internal/auth/transport"
	"github.com/nanmenkaimak/github-gist/internal/kafka"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo                     repository.Repository
	userGrpcTransport        *transport.UserGrpcTransport
	jwtSecretKey             string
	userVerificationProducer *kafka.Producer
	dbRedis                  *redis.Client
}

func NewAuthService(repo repository.Repository, authConfig config.Auth, userVerificationProducer *kafka.Producer,
	dbRedis *redis.Client, userGrpcTransport *transport.UserGrpcTransport) UseCase {
	return &Service{
		repo:                     repo,
		jwtSecretKey:             authConfig.JwtSecretKey,
		userVerificationProducer: userVerificationProducer,
		dbRedis:                  dbRedis,
		userGrpcTransport:        userGrpcTransport,
	}
}
