package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type router struct {
	logger *zap.SugaredLogger
}

func NewRouter(logger *zap.SugaredLogger) *router {
	return &router{
		logger: logger,
	}
}

func (s *router) GetHandler(eh *EndpointHandler) http.Handler {
	r := gin.Default()

	r.NoRoute(func(ctx *gin.Context) { // check for 404
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Page not found",
		})
	})

	auth := r.Group("/api/auth/v1")
	{
		auth.POST("/login", eh.Login)
		auth.POST("/renew-token", eh.RenewToken)
	}

	return r
}
