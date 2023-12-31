package http

import (
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
)

// swagger:route POST /v1/ Gists gist_create
//
// # Create Gist
//
// # Create Gist
//
// Consumes:
// - application/json
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: CreateGistRequest
//			in: body
//			required: true
//			type: CreateGistRequest
//
//		Security:
//		  Bearer:
//	 Responses:
//		  201: CreateGistResponse
//		  401:
//	   400:
func (h *EndpointHandler) CreateGist(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}

	var request dto.CreateGistRequest

	if err := ctx.BindJSON(&request); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	var gistFiles []entity.File

	for _, file := range request.FilesRequest {
		fileGist := entity.File{
			Name: file.Name,
			Code: file.Code,
		}
		gistFiles = append(gistFiles, fileGist)
	}

	newGist := entity.GistRequest{
		Gist: entity.Gist{
			Name:        request.GistRequest.Name,
			Description: request.GistRequest.Description,
			Visible:     request.GistRequest.Visible,
		},
		Commit: entity.Commit{
			Comment: request.CommitRequest.Comment,
		},
		Files: gistFiles,
	}

	newGist.Gist.UserID = userID.ID

	gistID, err := h.gistService.CreateGist(ctx.Request.Context(), newGist)
	if err != nil {
		h.logger.Errorf("failed to CreateGist err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, gistID)
}

// swagger:route GET /v1/discover Gists get_all_gists
//
// # Get All Gists
//
// # Get All Gists
//
//	Produces:
//	- application/json
//
//			Parameters:
//	      + name: sort
//	        in: query
//	        description: created_at, updated_at
//	        required: false
//	      + name: direction
//	        in: query
//	        description: asc, desc
//	        required: false
//
//	Schemes: http, https
//	Responses:
//	  200: []GistRequest
//	  400:
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

// swagger:route GET /v1/{username}/{gist_id} Gists get_gist_by_id
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
//		  + name: username
//			in: path
//		  + name: gist_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  200: GistRequest
//		  401:
//	      400:
func (h *EndpointHandler) GetGistByID(ctx *gin.Context) {
	var currentUserID uuid.UUID
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Warnf("cannot find user in context")
		currentUserID = uuid.Nil
	} else {
		currentUserID = userID.ID
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
		UserID:   currentUserID,
	}

	gist, err := h.gistService.GetGistByID(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to GetGistByID err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gist)
}

// swagger:route GET /v1/{username}/gists Gists get_all_gists_of_user
//
// # Get All Gists Of User
//
// # If it is user's account, they can see secret gists
//
// Produces:
// -application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: username
//			in: path
//	      + name: q
//	        in: query
//	        description: search by name
//	        required: false
//
//		Security:
//		  Bearer:
//	 Responses:
//		  200: []GistRequest
//		  401:
//	   400:
func (h *EndpointHandler) GetAllGistsOfUser(ctx *gin.Context) {
	var currentUserID uuid.UUID
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Warnf("cannot find user in context")
		currentUserID = uuid.Nil
	} else {
		currentUserID = userID.ID
	}
	username := ctx.Param("username")
	q := ctx.Query("q")

	request := gist.GetGistRequest{
		Username:  username,
		UserID:    currentUserID,
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

// swagger:route GET /v1/{username}/secret Gists get_all_secret_gists
//
// # Get All Secret Gists
//
// # Get All Secret Gists, if it is user's account
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
//		  401:
//	   400:
func (h *EndpointHandler) GetAllSecretGists(ctx *gin.Context) {
	var currentUserID uuid.UUID
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Warnf("cannot find user in context")
		currentUserID = uuid.Nil
	} else {
		currentUserID = userID.ID
	}
	username := ctx.Param("username")
	request := gist.GetGistRequest{
		Username:   username,
		UserID:     currentUserID,
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

// swagger:route GET /v1/{username}/public Gists get_all_public_gists
//
// # Get All Public Gists
//
// # Get All Public Gists Of Users
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
//		  401:
//	   400:
func (h *EndpointHandler) GetAllPublicGists(ctx *gin.Context) {
	var currentUserID uuid.UUID
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Warnf("cannot find user in context")
		currentUserID = uuid.Nil
	} else {
		currentUserID = userID.ID
	}
	username := ctx.Param("username")
	request := gist.GetGistRequest{
		Username:   username,
		UserID:     currentUserID,
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

// swagger:route PUT /v1/{username}/{gist_id} Gists update_gist_by_id
//
// # Update Gist By ID
//
// # Update Gist By ID, if it is user's account
//
// Consumes:
// - application/json
//
//		Schemes: http, https
//		Parameters:
//		  + name: CreateGistRequest
//			in: body
//			required: true
//			type: CreateGistRequest
//		  + name: username
//			in: path
//		  + name: gist_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  204:
//		  401:
//	   400:
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

	var request dto.CreateGistRequest

	if err := ctx.BindJSON(&request); err != nil {
		h.logger.Errorf("failed to unmarshall body err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	var gistFiles []entity.File

	for _, file := range request.FilesRequest {
		fileGist := entity.File{
			Name: file.Name,
			Code: file.Code,
		}
		gistFiles = append(gistFiles, fileGist)
	}

	updatedGist := entity.GistRequest{
		Gist: entity.Gist{
			ID:          gistID,
			Name:        request.GistRequest.Name,
			Description: request.GistRequest.Description,
			Visible:     request.GistRequest.Visible,
		},
		Commit: entity.Commit{
			Comment: request.CommitRequest.Comment,
		},
		Files: gistFiles,
	}

	updatedGist.Gist.UserID = userID.ID

	updateGistRequest := gist.UpdateGistRequest{
		Username: username,
		Gist:     updatedGist,
		UserID:   userID.ID,
	}

	err = h.gistService.UpdateGistByID(ctx.Request.Context(), updateGistRequest)
	if err != nil {
		h.logger.Errorf("failed to UpdateGistByID err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// swagger:route DELETE /v1/{username}/{gist_id} Gists delete_gist_by_id
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
//		  + name: username
//			in: path
//		  + name: gist_id
//			in: path
//
//		Security:
//		  Bearer:
//	 Responses:
//		  204:
//		  401:
//	   400:
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
