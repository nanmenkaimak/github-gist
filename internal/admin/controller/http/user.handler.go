package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanmenkaimak/github-gist/internal/admin/admin"
	"github.com/nanmenkaimak/github-gist/internal/admin/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
)

func (h *EndpointHandler) UpdateUserByUsername(ctx echo.Context) error {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		return ctx.NoContent(http.StatusUnauthorized)
	}
	username := ctx.Param("username")

	var updatedUser entity.User

	if err := ctx.Bind(&updatedUser); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	updatedUser.Username = username

	request := admin.UpdateUserRequest{
		UpdatedUser: updatedUser,
		UserID:      userID.ID,
	}

	err = h.adminService.UpdateUserByUsername(ctx.Request().Context(), request)
	if err != nil {
		h.logger.Errorf("failed to UpdateUser err: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *EndpointHandler) GetAllUsers(ctx echo.Context) error {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		return ctx.NoContent(http.StatusUnauthorized)
	}

	request := admin.GetUserRequest{
		UserID: userID.ID,
	}

	allUsers, err := h.adminService.GetAllUsers(ctx.Request().Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetAllUsers err: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, allUsers)
}

func (h *EndpointHandler) GetUserByUsername(ctx echo.Context) error {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		return ctx.NoContent(http.StatusUnauthorized)
	}
	username := ctx.Param("username")

	request := admin.GetUserRequest{
		UserID:              userID.ID,
		UpdatedUserUsername: username,
	}
	user, err := h.adminService.GetUserByUsername(ctx.Request().Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetAllUsers err: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, user)
}
