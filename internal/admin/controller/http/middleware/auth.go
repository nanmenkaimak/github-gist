package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nanmenkaimak/github-gist/internal/admin/auth"
	"go.uber.org/zap"
)

const (
	AuthorizationHeaderKey = "Authorization"
)

type JwtV1 struct {
	authService auth.UseCase
	logger      *zap.SugaredLogger
}

func NewJwtV1Middleware(authService auth.UseCase, logger *zap.SugaredLogger) *JwtV1 {
	return &JwtV1{
		authService: authService,
		logger:      logger,
	}
}

func (j *JwtV1) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearer := c.Request().Header.Get(AuthorizationHeaderKey)
		if bearer == "" {
			j.logger.Warn("'Authorization' key missing from headers")
			return c.NoContent(http.StatusUnauthorized)
		}
		var jwtToken string
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			jwtToken = bearer[7:]
		} else {
			j.logger.Warn(fmt.Sprintf(
				"failed to get token from header invalidToken: %s",
				c.Request().Header.Get(AuthorizationHeaderKey),
			))
			return c.NoContent(http.StatusUnauthorized)
		}

		contextUser, err := j.authService.GetJWTUser(jwtToken)
		if err != nil {
			j.logger.Errorf("failed to GetJwtUser err: %v", err)
			return c.NoContent(http.StatusUnauthorized)
		}

		if contextUser.RoleID != 2 {
			j.logger.Errorf("only admin have access err: %v", err)
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set("id", contextUser)
		return next(c)
	}
}

func GetContextUser(c echo.Context) (*auth.ContextUser, error) {
	id := c.Get("id")
	if id == nil {
		return nil, errors.New("cannot get ID")
	}
	return id.(*auth.ContextUser), nil
}
