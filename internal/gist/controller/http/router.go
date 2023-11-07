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
		gist.GET("/:username/:gist_id", eh.GetGistByID)
		gist.GET("/:username/gists", eh.GetAllGistsOfUser)
		gist.GET("/:username/secret", eh.GetAllSecretGists)
		gist.GET("/:username/public", eh.GetAllPublicGists)
		gist.PUT("/:username/:gist_id/edit", eh.UpdateGistByID)
		gist.DELETE("/:username/:gist_id/edit", eh.DeleteGistByID)

		gist.POST("/:username/:gist_id/star", eh.StarGist)
		gist.GET("/:username/starred", eh.GetStaredGists)
		gist.DELETE("/:username/:gist_id/star", eh.DeleteStar)

		gist.POST("/:username/:gist_id/fork", eh.ForkGist)
		gist.GET("/:username/forked", eh.GetForkedGists)

		gist.POST("/:username/:gist_id/comment", eh.CreateComment)
		gist.GET("/:username/:gist_id/comment", eh.GetCommentsOfGist)
		gist.DELETE("/:username/:gist_id/comment/:comment_id", eh.DeleteComment)
		gist.PATCH("/:username/:gist_id/comment/:comment_id", eh.UpdateComment)

		gist.POST("/:username/follow", eh.FollowUser)
		gist.POST("/:username/unfollow", eh.UnfollowUser)
	}

	return r
}
