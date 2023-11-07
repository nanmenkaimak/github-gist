package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/user/entity"
)

func (r *Repo) FollowUser(follower entity.Follower) error {
	var alreadyFollowed entity.Follower
	err := r.replica.Db.Where("follower_id = ? and following_id = ?", follower.FollowerID, follower.FollowingID).
		Find(&alreadyFollowed).Error
	if err != nil {
		return err
	}

	if alreadyFollowed.ID != uuid.Nil {
		return fmt.Errorf("you already followed user")
	}

	err = r.main.Db.Create(&follower).Error
	return err
}

func (r *Repo) UnfollowUser(follower entity.Follower) error {
	err := r.main.Db.Where("follower_id = ? and following_id = ?", follower.FollowerID, follower.FollowingID).
		Delete(&entity.Follower{}).Error
	return err
}
