package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/auth/entitiy"
	"net/http"
)

func (f *EndpointHandler) Register(ctx *gin.Context) {
	var request entitiy.RegisterUserRequest
	if err := ctx.BindJSON(&request); err != nil {
		f.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	// make register user
	userID, err := f.authService.RegisterUser(ctx.Request.Context(), request)
	if err != nil {
		f.logger.Errorf("failed to Register err: %v", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusCreated, userID)
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
	ctx.Status(http.StatusNoContent)
}

func (f *EndpointHandler) UpdateUser(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		f.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")

	var updatedUser entitiy.RegisterUserRequest

	if err := ctx.BindJSON(&updatedUser); err != nil {
		f.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	updatedUser.Username = username
	updatedUser.ID = userID.ID

	err = f.authService.UpdateUser(ctx.Request.Context(), updatedUser)
	if err != nil {
		f.logger.Errorf("failed to UpdateUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (f *EndpointHandler) ResetCode(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		f.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")

	request := auth.ResetCodeRequest{
		Username: username,
		UserID:   userID.ID,
	}

	err = f.authService.ResetCode(ctx.Request.Context(), request)
	if err != nil {
		f.logger.Errorf("failed to ResetCode err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (f *EndpointHandler) ResetPassword(ctx *gin.Context) {
	username := ctx.Param("username")

	var request auth.UpdatePasswordRequest

	if err := ctx.BindJSON(&request); err != nil {
		f.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	request.Username = username

	err := f.authService.ResetPassword(ctx.Request.Context(), request)
	if err != nil {
		f.logger.Errorf("failed to ResetPassword err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}
