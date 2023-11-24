package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nanmenkaimak/github-gist/internal/admin/admin"
	"github.com/nanmenkaimak/github-gist/internal/admin/controller/http/middleware"
)

// swagger:route GET /v1/gist/ Gists get_all_gists
//
// # Get All Gists
//
// Get All Gists
//
//	Produces:
//	- application/json
//
//		Security:
//		  Bearer:
//
//	Schemes: http, https
//	Responses:
//	  200: []GistRequest
//	  400:
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

// swagger:route GET /v1/gist/{gist_id} Gists get_gist_by_id
//
// # Get Gist By ID
//
// # Get Gist By ID of user
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: gist_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  200: GistRequest
//		  401:
//	      400:
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

// swagger:route DELETE /v1/gist/{gist_id} Gists delete_gist_by_id
//
// # Delete Gist By ID
//
// # Delete Gist By ID, if it user's account
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: gist_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  204:
//		  401:
//	   400:
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
