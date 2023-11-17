package http

import (
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"go.uber.org/zap"
)

type EndpointHandler struct {
	gistService gist.UseCase
	logger      *zap.SugaredLogger
}

func NewEndpointHandler(gistService gist.UseCase, logger *zap.SugaredLogger) *EndpointHandler {
	return &EndpointHandler{
		gistService: gistService,
		logger:      logger,
	}
}
