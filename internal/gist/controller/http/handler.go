package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"go.uber.org/zap"
	"net/http"
)

type EndpointHandler struct {
	gistService gist.UseCase
	logger      *zap.SugaredLogger
}

func NewEndpointHandler(gistService gist.UseCase, logger *zap.SugaredLogger) *EndpointHandler {
	return &EndpointHandler{
		gistService: gistService,
		logger:      logger,
	}
}

func (h *EndpointHandler) CreateGist(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var request entity.GistRequest

	if err := ctx.BindJSON(&request); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	request.Gist.UserID = userID.ID

	gistID, err := h.gistService.CreateGist(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to CreateGist err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, gistID)
}

func (h *EndpointHandler) GetAllGists(ctx *gin.Context) {
	allGists, err := h.gistService.GetAllGists(ctx.Request.Context())
	if err != nil {
		h.logger.Errorf("failed to GetAllGists err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, allGists)
}

func (h *EndpointHandler) GetGistByID(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}

	gistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return
	}
	username := ctx.Param("username")

	request := gist.GetGistRequest{
		GistID:   gistID,
		Username: username,
		UserID:   userID.ID,
	}

	gist, err := h.gistService.GetGistByID(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetGistByID err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gist)
}

func (h *EndpointHandler) GetAllGistsOfUser(ctx *gin.Context) {
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

	gists, err := h.gistService.GetAllGistsOfUser(ctx, request)
	if err != nil {
		h.logger.Errorf("failed to GetAllGistsOfUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, gists)
}

func (h *EndpointHandler) UpdateGistByID(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	gistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return
	}
	username := ctx.Param("username")

	var updatedGist entity.GistRequest

	if err := ctx.BindJSON(&updatedGist); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	updatedGist.Gist.ID = gistID

	request := gist.UpdateGistRequest{
		Username: username,
		Gist:     updatedGist,
		UserID:   userID.ID,
	}

	err = h.gistService.UpdateGistByID(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to UpdateGistByID err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *EndpointHandler) DeleteGistByID(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	gistID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return
	}
	username := ctx.Param("username")

	request := gist.GetGistRequest{
		Username: username,
		GistID:   gistID,
		UserID:   userID.ID,
	}

	err = h.gistService.DeleteGistByID(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to DeleteGistByID err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}
