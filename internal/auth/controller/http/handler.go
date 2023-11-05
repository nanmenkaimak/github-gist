package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
	"go.uber.org/zap"
	"net/http"
	"unicode"
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

func (f *EndpointHandler) Register(ctx *gin.Context) {
	// make register user
	err := f.authService.RegisterUser(ctx.Request.Context())
	if err != nil {
		f.logger.Errorf("failed to Register err: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
}

func (f *EndpointHandler) ConfirmUser(ctx *gin.Context) {
	request := struct {
		Code  string `json:"code"`
		Email string `json:"email"`
	}{}
	if err := ctx.BindJSON(&request); err != nil {
		f.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	confirmRequest := auth.ConfirmUserRequest{
		Code:  request.Code,
		Email: request.Email,
	}
	err := f.authService.ConfirmUser(ctx.Request.Context(), confirmRequest)
	if err != nil {
		f.logger.Errorf("failed to ConfirmUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, "norm")
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

	err := validPassword(request.Password)
	if err != nil {
		f.logger.Errorf("validPassword err: %v", err)
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

func validPassword(s string) error {
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
		"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("password must have at least one %s character", name)
	}
	return nil
}
