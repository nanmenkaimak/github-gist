package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/gist/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"github.com/nanmenkaimak/github-gist/internal/gist/gist"
	"net/http"
)

func (h *EndpointHandler) FollowUser(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")

	request := gist.FollowRequest{
		FollowerID: userID.ID,
		Username:   username,
	}

	err = h.gistService.FollowUser(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to FollowUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *EndpointHandler) UnfollowUser(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		h.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")

	request := gist.FollowRequest{
		FollowerID: userID.ID,
		Username:   username,
	}

	err = h.gistService.UnfollowUser(ctx.Request.Context(), request)
	if err != nil {
		h.logger.Errorf("failed to UnfollowUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (h *EndpointHandler) UserInfo(ctx *gin.Context) {
	username := ctx.Param("username")
	tab := ctx.Query("tab")
	response := &[]entity.UserResponse{}
	var err error
	if tab == "follower" {
		response, err = h.gistService.GetAllFollowers(ctx.Request.Context(), username)
		if err != nil {
			h.logger.Errorf("failed to GetAllFollowers err: %v", err)
		}
	} else if tab == "following" {
		response, err = h.gistService.GetAllFollowings(ctx.Request.Context(), username)
		if err != nil {
			h.logger.Errorf("failed to GetAllFollowings err: %v", err)
		}
	} else {
		responseInfo, err := h.gistService.GetUserInfo(ctx.Request.Context(), username)
		if err != nil {
			h.logger.Errorf("failed to GetUserInfo err: %v", err)
		}
		*response = append(*response, responseInfo)
	}
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
