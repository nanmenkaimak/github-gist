package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/http/middleware"
	"go.uber.org/zap"
	"net/http"
)

type router struct {
	logger         *zap.SugaredLogger
	authMiddleware *middleware.JwtV1
}

func NewRouter(logger *zap.SugaredLogger, authMiddleware *middleware.JwtV1) *router {
	return &router{
		logger:         logger,
		authMiddleware: authMiddleware,
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
		auth.POST("/register", eh.Register)
		auth.POST("/confirm-user", eh.ConfirmUser)
		auth.POST("/:username/reset-code")
		auth.PATCH("/:username/reset-password")
	}

	auth.Use(s.authMiddleware.Auth())
	{
		auth.PUT("/:username/update", eh.UpdateUser)
	}

	return r
}
