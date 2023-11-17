package http

import (
	"github.com/nanmenkaimak/github-gist/internal/admin/admin"
	"go.uber.org/zap"
)

type EndpointHandler struct {
	adminService admin.UseCase
	logger       *zap.SugaredLogger
}

func NewEndpointHandler(adminService admin.UseCase, logger *zap.SugaredLogger) *EndpointHandler {
	return &EndpointHandler{
		adminService: adminService,
		logger:       logger,
	}
}
