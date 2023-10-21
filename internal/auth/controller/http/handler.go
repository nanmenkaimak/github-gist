package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
	"go.uber.org/zap"
	"net/http"
)

type EndpointHandler struct {
	authService auth.UseCase
	logger      *zap.SugaredLogger
}

func NewEndpointHandler(authService auth.UseCase, logger *zap.SugaredLogger) *EndpointHandler {
	return &EndpointHandler{
		authService: authService,
		logger:      logger,
	}
}

func (f *EndpointHandler) Login(ctx *gin.Context) {
	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := ctx.BindJSON(&request); err != nil {
		f.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	generateTokenRequest := auth.GenerateTokenRequest{
		Username: request.Username,
		Password: request.Password,
	}

	userToken, err := f.authService.GenerateToken(ctx.Request.Context(), generateTokenRequest)
	if err != nil {
		f.logger.Errorf("failed to GenerateToken err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	response := struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Token:        userToken.Token,
		RefreshToken: userToken.RefreshToken,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (f *EndpointHandler) RenewToken(ctx *gin.Context) {
	request := struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := ctx.BindJSON(&request); err != nil {
		f.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	userToken, err := f.authService.RenewToken(ctx.Request.Context(), request.RefreshToken)
	if err != nil {
		f.logger.Errorf("failed to RenewToken err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: userToken.Token,
	}

	ctx.JSON(http.StatusCreated, response)
}
