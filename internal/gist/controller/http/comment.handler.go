package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"net/http"
)

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

func (h *EndpointHandler) DeleteComment(ctx *gin.Context) {
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
	commentID, err := uuid.Parse(ctx.Param("comment_id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return
	}
	username := ctx.Param("username")

	request := gist.DeleteRequest{
		GistID:    gistID,
		Username:  username,
		UserID:    userID.ID,
		CommentID: commentID,
	}

	err = h.gistService.DeleteComment(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to DeleteComment err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (h *EndpointHandler) UpdateComment(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	commentID, err := uuid.Parse(ctx.Param("comment_id"))
	if err != nil {
		h.logger.Errorf("parsing value from url err: %v", err)
		return
	}
	username := ctx.Param("username")

	var updatedComment struct {
		Text string `json:"text"`
	}

	if err := ctx.BindJSON(&updatedComment); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	request := gist.UpdateCommentRequest{
		CommentID: commentID,
		Username:  username,
		UserID:    userID.ID,
		Text:      updatedComment.Text,
	}

	err = h.gistService.UpdateComment(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to UpdateComment err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusNoContent)
}
