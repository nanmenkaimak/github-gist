package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
)

// swagger:route POST /v1/{username}/{gist_id}/fork Fork fork_gist
//
// # Fork Gist
//
// # Fork Gist
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//	      + name: username
//			in: path
//		  + name: gist_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  201:
//		  401:
//	      400:
func (h *EndpointHandler) ForkGist(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	gistID, err := uuid.Parse(ctx.Param("gist_id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return
	}
	username := ctx.Param("username")

	request := gist.ForkRequest{
		UserID:   userID.ID,
		Username: username,
		GistID:   gistID,
	}

	response, err := h.gistService.ForkGist(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to ForkGist err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// swagger:route GET /v1/{username}/forked Fork get_forked_gists
//
// # Get Forked Gists
//
// # Get Forked Gists
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: username
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  200: []GistRequest
//	      400:
func (h *EndpointHandler) GetForkedGists(ctx *gin.Context) {
	var currentUserID uuid.UUID
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Warn("cannot find user in context")
		currentUserID = uuid.Nil
	} else {
		currentUserID = userID.ID
	}

	username := ctx.Param("username")
	request := gist.GetGistRequest{
		Username: username,
		UserID:   currentUserID,
	}

	forkedGists, err := h.gistService.GetForkedGists(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to ForkGist err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, forkedGists)
}
