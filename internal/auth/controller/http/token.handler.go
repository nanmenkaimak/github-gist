package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
)

// swagger:route POST /v1/login Token login
//
// # Login User
//
// # Login User
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: GenerateTokenRequest
//			in: body
//			required: true
//			type: GenerateTokenRequest
//
//	 Responses:
//		  201: GenerateTokenResponse
//	   400:
func (f *EndpointHandler) Login(ctx *gin.Context) {
	var generateTokenRequest auth.GenerateTokenRequest

	if err := ctx.BindJSON(&generateTokenRequest); err != nil {
		f.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err := validPassword(generateTokenRequest.Password)
	if err != nil {
		f.logger.Errorf("validPassword err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	userToken, err := f.authService.GenerateToken(ctx.Request.Context(), generateTokenRequest)
	if err != nil {
		f.logger.Errorf("failed to GenerateToken err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	response := auth.GenerateTokenResponse{
		Token:        userToken.Token,
		RefreshToken: userToken.RefreshToken,
	}

	ctx.JSON(http.StatusCreated, response)
}

// swagger:route POST /v1/renew-token Token renew_token
//
// # Renew Token
//
// # Renew token of user
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: RenewTokenRequest
//			in: body
//			required: true
//			type: RenewTokenRequest
//
//	 Responses:
//		  201: RenewTokenResponse
//	   400:
func (f *EndpointHandler) RenewToken(ctx *gin.Context) {
	var request auth.RenewTokenRequest

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

	response := auth.RenewTokenResponse{
		Token: userToken.Token,
	}

	ctx.JSON(http.StatusCreated, response)
}
