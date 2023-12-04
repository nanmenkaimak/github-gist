package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
)

// swagger:route POST /v1/register User register
//
// # Register User
//
// # Register User
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: RegisterUserRequest
//			in: body
//			required: true
//			type: RegisterUserRequest
//
//	 Responses:
//		  201: RegisterUserResponse
//	   400:
func (f *EndpointHandler) Register(ctx *gin.Context) {
	var request entity.RegisterUserRequest
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

// swagger:route POST /v1/confirm-user User confirm_user
//
// # Confirm User
//
// # Confirm User
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: ConfirmUserRequest
//			in: body
//			required: true
//			type: ConfirmUserRequest
//
//	 Responses:
//		  204:
//	   400:
func (f *EndpointHandler) ConfirmUser(ctx *gin.Context) {
	var request auth.ConfirmUserRequest
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

// swagger:route PUT /v1/{username}/update User update_user
//
// # Update User
//
// # Update User
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: RegisterUserRequest
//			in: body
//			required: true
//			type: RegisterUserRequest
//		  + name: username
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  204:
//		  401:
//	   400:
func (f *EndpointHandler) UpdateUser(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		f.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")

	var updatedUser entity.RegisterUserRequest

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

// swagger:route POST /v1/{username}/reset-code User reset_code
//
// # Send Reset Code
//
// # Send Reset Code For Updating Password
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: username
//			in: path
//
//	 Responses:
//		  204:
//	   400:
func (f *EndpointHandler) ResetCode(ctx *gin.Context) {
	username := ctx.Param("username")

	request := auth.ResetCodeRequest{
		Username: username,
	}

	err := f.authService.ResetCode(ctx.Request.Context(), request)
	if err != nil {
		f.logger.Errorf("failed to ResetCode err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// swagger:route PATCH /v1/{username}/reset-password User reset_password
//
// # Update Password
//
// # Update Password
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: UpdatePasswordRequest
//			in: body
//			required: true
//			type: UpdatePasswordRequest
//		  + name: username
//			in: path
//
//	 Responses:
//		  204:
//	   400:
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
