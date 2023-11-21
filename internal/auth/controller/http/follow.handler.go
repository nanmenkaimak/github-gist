package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nanmenkaimak/github-gist/internal/auth/auth"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/http/middleware"
	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
	"net/http"
)

func (f *EndpointHandler) FollowUser(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		f.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")

	request := auth.FollowRequest{
		FollowerID: userID.ID,
		Username:   username,
	}

	err = f.authService.FollowUser(ctx.Request.Context(), request)
	if err != nil {
		f.logger.Errorf("failed to FollowUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (f *EndpointHandler) UnfollowUser(ctx *gin.Context) {
	userID, err := middleware.GetContextUser(ctx)
	if err != nil {
		f.logger.Errorf("cannot find user in context")
		ctx.Status(http.StatusUnauthorized)
		return
	}
	username := ctx.Param("username")

	request := auth.FollowRequest{
		FollowerID: userID.ID,
		Username:   username,
	}

	err = f.authService.UnfollowUser(ctx.Request.Context(), request)
	if err != nil {
		f.logger.Errorf("failed to UnfollowUser err: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (f *EndpointHandler) UserInfo(ctx *gin.Context) {
	username := ctx.Param("username")
	tab := ctx.Query("tab")
	response := &[]entity.RegisterUserRequest{}
	var err error
	if tab == "follower" {
		response, err = f.authService.GetAllFollowers(ctx.Request.Context(), username)
		if err != nil {
			f.logger.Errorf("failed to GetAllFollowers err: %v", err)
		}
	} else if tab == "following" {
		response, err = f.authService.GetAllFollowings(ctx.Request.Context(), username)
		if err != nil {
			f.logger.Errorf("failed to GetAllFollowings err: %v", err)
		}
	} else {
		responseInfo, err := f.authService.GetUserInfo(ctx.Request.Context(), username)
		if err != nil {
			f.logger.Errorf("failed to GetUserInfo err: %v", err)
		}
		*response = append(*response, responseInfo)
	}
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
