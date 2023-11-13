package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nanmenkaimak/github-gist/internal/admin/admin"
	"github.com/nanmenkaimak/github-gist/internal/admin/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
	"go.uber.org/zap"
	"net/http"
)

type EndpointHandler struct {
	adminService admin.UseCase
	logger       *zap.SugaredLogger
}

func NewEndpointHandler(adminService admin.UseCase, logger *zap.SugaredLogger) *EndpointHandler {
	return &EndpointHandler{
		adminService: adminService,
		logger:       logger,
	}
}

func (h *EndpointHandler) GetAllGists(ctx echo.Context) error {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		return ctx.NoContent(http.StatusUnauthorized)
	}
	direction := ctx.QueryParam("direction")
	sort := ctx.QueryParam("sort")
	request := admin.GetAllGistsRequest{
		Direction: direction,
		Sort:      sort,
		UserID:    userID.ID,
	}

	allGists, err := h.adminService.GetAllGists(ctx.Request().Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetAllGists err: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, allGists)
}

func (h *EndpointHandler) GetGistByID(ctx echo.Context) error {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		return ctx.NoContent(http.StatusUnauthorized)
	}
	gistID, err := uuid.Parse(ctx.Param("gist_id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return nil
	}

	request := admin.GetGistRequest{
		GistID: gistID,
		UserID: userID.ID,
	}

	gist, err := h.adminService.GetGistByID(ctx.Request().Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetGistByID err: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, gist)
}

func (h *EndpointHandler) DeleteGistByID(ctx echo.Context) error {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		return ctx.NoContent(http.StatusUnauthorized)
	}
	gistID, err := uuid.Parse(ctx.Param("gist_id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return nil
	}

	request := admin.GetGistRequest{
		GistID: gistID,
		UserID: userID.ID,
	}

	err = h.adminService.DeleteGistByID(ctx.Request().Context(), request)
	if err != nil {
		h.logger.Errorf("failed to DeleteGistByID err: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.NoContent(http.StatusNoContent)
}

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
