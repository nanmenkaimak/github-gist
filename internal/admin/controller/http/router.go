package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/nanmenkaimak/github-gist/internal/admin/controller/http/middleware"
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
	r := echo.New()

	r.Use(echo_middleware.Logger())
	r.Use(echo_middleware.Recover())

	admin := r.Group("/api/admin/v1")
	admin.Use(s.authMiddleware.Auth)

	gist := admin.Group("/gist")
	gist.GET("/", eh.GetAllGists)
	gist.GET("/:gist_id", eh.GetGistByID)
	gist.DELETE("/:gist_id", eh.DeleteGistByID)

	user := admin.Group("/user")
	user.PUT("/:username", eh.UpdateUserByUsername)
	user.GET("/", eh.GetAllUsers)
	user.GET("/:username", eh.GetUserByUsername)

	return r
}
