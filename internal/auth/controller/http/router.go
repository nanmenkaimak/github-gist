package http

import (
	"github.com/gin-contrib/cors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/http/middleware"
	"go.uber.org/zap"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8082"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type", "Accept-Encoding"},
	}))

	auth := r.Group("/api/auth/v1")
	{
		auth.POST("/login", eh.Login)
		auth.POST("/renew-token", eh.RenewToken)
		auth.POST("/register", eh.Register)
		auth.POST("/confirm-user", eh.ConfirmUser)
		auth.GET("/:username", eh.UserInfo)
	}

	auth.Use(s.authMiddleware.Auth())
	{
		auth.PUT("/:username/update", eh.UpdateUser)
		auth.POST("/:username/reset-code", eh.ResetCode)
		auth.PATCH("/:username/reset-password", eh.ResetPassword)
		auth.POST("/:username/follow", eh.FollowUser)
		auth.POST("/:username/unfollow", eh.UnfollowUser)
	}

	return r
}
