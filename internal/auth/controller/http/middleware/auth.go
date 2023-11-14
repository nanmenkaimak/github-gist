package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
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

func (j *JwtV1) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader(AuthorizationHeaderKey)
		if bearer == "" {
			j.logger.Warn("'Authorization' key missing from headers")
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		var jwtToken string
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			jwtToken = bearer[7:]
		} else {
			j.logger.Warn(fmt.Sprintf(
				"failed to get token from header invalidToken: %s",
				c.GetHeader(AuthorizationHeaderKey),
			))
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		contextUser, err := j.authService.GetJWTUser(jwtToken)
		if err != nil {
			j.logger.Errorf("failed to GetJwtUser err: %v", err)
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.Set("id", contextUser)
		c.Next()
	}
}

func GetContextUser(ctx *gin.Context) (*auth.ContextUser, error) {
	id, ok := ctx.Get("id")
	if !ok {
		return nil, errors.New("cannot get ID")
	}
	return id.(*auth.ContextUser), nil
}
