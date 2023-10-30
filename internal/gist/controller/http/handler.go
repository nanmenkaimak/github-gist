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
	direction := ctx.Query("direction")
	sort := ctx.Query("sort")
	request := gist.GetAllGistsRequest{
		Direction: direction,
		Sort:      sort,
	}
	allGists, err := h.gistService.GetAllGists(ctx.Request.Context(), request)
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

	gistID, err := uuid.Parse(ctx.Param("gist_id"))
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
	q := ctx.Query("q")

	request := gist.GetGistRequest{
		Username:  username,
		UserID:    userID.ID,
		Searching: q,
	}

	gists, err := h.gistService.GetAllGistsOfUser(ctx, request)
	if err != nil {
		h.logger.Errorf("failed to GetAllGistsOfUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, gists)
}

func (h *EndpointHandler) GetAllSecretGists(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")
	request := gist.GetGistRequest{
		Username:   username,
		UserID:     userID.ID,
		Visibility: false,
	}

	gists, err := h.gistService.GetGistsByVisibility(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetAllSecretGists err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, gists)
}

func (h *EndpointHandler) GetAllPublicGists(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")
	request := gist.GetGistRequest{
		Username:   username,
		UserID:     userID.ID,
		Visibility: true,
	}

	gists, err := h.gistService.GetGistsByVisibility(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetAllSecretGists err: %v", err)
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
	gistID, err := uuid.Parse(ctx.Param("gist_id"))
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
	gistID, err := uuid.Parse(ctx.Param("gist_id"))
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

func (h *EndpointHandler) CreateComment(ctx *gin.Context) {
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
	var request entity.Comment

	if err := ctx.BindJSON(&request); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	request.GistID = gistID
	request.UserID = userID.ID

	err = h.gistService.CreateComment(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to CreateComment err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *EndpointHandler) GetCommentsOfGist(ctx *gin.Context) {
	gistID, err := uuid.Parse(ctx.Param("gist_id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return
	}

	request := gist.GetGistRequest{
		GistID: gistID,
	}

	comments, err := h.gistService.GetCommentsOfGist(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetCommentsOfGist err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
