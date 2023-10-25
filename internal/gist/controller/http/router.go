package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
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

	gist := r.Group("/api/gist/v1")
	{
		gist.Use(s.authMiddleware.Auth())
		gist.POST("/create-gist", eh.CreateGist)
		gist.GET("/all-gists", eh.GetAllGists)
		gist.GET("/:username/:id", eh.GetGistByID)
		gist.GET("/:username/gists", eh.GetAllGistsOfUser)
		gist.GET("/:username/secret", eh.GetAllSecretGists)
		gist.GET("/:username/public", eh.GetAllPublicGists)
		// post searching gists
		gist.PUT("/:username/:id/edit", eh.UpdateGistByID)
		gist.DELETE("/:username/:id/edit", eh.DeleteGistByID)
	}

	return r
}
