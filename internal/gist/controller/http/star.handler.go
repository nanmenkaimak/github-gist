package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"net/http"
)

func (h *EndpointHandler) StarGist(ctx *gin.Context) {
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

	request := entity.Star{
		UserID: userID.ID,
		GistID: gistID,
	}

	err = h.gistService.StarGist(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to StarGist err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *EndpointHandler) GetStaredGists(ctx *gin.Context) {
	username := ctx.Param("username")
	request := gist.OtherGistRequest{
		Username: username,
	}

	gists, err := h.gistService.GetStaredGists(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetStaredGists err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gists)
}

func (h *EndpointHandler) DeleteStar(ctx *gin.Context) {
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

	request := gist.DeleteRequest{
		GistID:   gistID,
		Username: username,
		UserID:   userID.ID,
	}

	err = h.gistService.DeleteStar(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to DeleteStar err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusNoContent)
}
