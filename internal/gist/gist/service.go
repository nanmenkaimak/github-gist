package gist

import (
	"github.com/nanmenkaimak/github-gist/internal/gist/repository"
	"github.com/nanmenkaimak/github-gist/internal/gist/transport"
)

type Service struct {
	repo              repository.Repository
	userGrpcTransport *transport.UserGrpcTransport
}

func NewGistService(repo repository.Repository, userGrpcTransport *transport.UserGrpcTransport) UseCase {
	return &Service{
		repo:              repo,
		userGrpcTransport: userGrpcTransport,
	}
}
