package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"net/http"
)

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

func (h *EndpointHandler) GetForkedGists(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")
	request := gist.GetGistRequest{
		Username: username,
		UserID:   userID.ID,
	}

	forkedGists, err := h.gistService.GetForkedGists(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to ForkGist err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, forkedGists)
}
