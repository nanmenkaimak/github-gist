package http

import (
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
)

// swagger:route POST /v1/{username}/{gist_id}/comment Comment comment_create
//
// # Create Comment
//
// # Create Comment
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: Comment
//			in: body
//			required: true
//			type: Comment
//		  + name: username
//			in: path
//		  + name: gist_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  201:
//		  401:
//	   400:
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
	var request = struct {
		Text string `json:"text"`
	}{}

	if err := ctx.BindJSON(&request); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	commentBuilder := entity.NewCommentBuilder()
	newComment := commentBuilder.SetText(request.Text).SetGistID(gistID).SetUserID(userID.ID).Build()

	err = h.gistService.CreateComment(ctx.Request.Context(), *newComment)
	if err != nil {
		h.logger.Errorf("failed to CreateComment err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusCreated)
}

// swagger:route GET /v1/{username}/{gist_id}/comment Comment get_comments
//
// # Get Comments Of Gist
//
// # Get Comments Of Gist
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: username
//			in: path
//		  + name: gist_id
//			in: path
//
//	 Responses:
//		  200: []Comment
//	      400:
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

// swagger:route DELETE /v1/{username}/{gist_id}/comment/{comment_id} Comment delete_comment_by_id
//
// # Delete Comment By ID
//
// # Delete Comment By ID, if it user's comment
//
// Produces:
// -application/json
//
//			Schemes: http, https
//			Parameters:
//			  + name: username
//				in: path
//			  + name: gist_id
//				in: path
//	       + name: comment_id
//				in: path
//
//			Security:
//			  Bearer:
//		 Responses:
//			  204:
//			  401:
//		   400:
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

// swagger:route PATCH /v1/{username}/{gist_id}/comment/{comment_id} Comment update_comment_by_id
//
// # Update Comment By ID
//
// # Update Comment By ID, if it is user's comment
//
// Consumes:
// - application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: Comment
//			in: body
//			required: true
//			type: Comment
//		  + name: username
//			in: path
//		  + name: gist_id
//			in: path
//		  + name: comment_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  204:
//		  401:
//	   400:
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
