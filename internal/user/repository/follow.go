package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/user/entity"
	"gorm.io/gorm"
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

func (r *Repo) GetAllFollowers(userID string) ([]entity.User, error) {
	var followers []entity.User
	err := r.replica.Db.Transaction(func(tx *gorm.DB) error {
		var allFollowersID []entity.Follower
		res := tx.Where("following_id = ?", userID).Find(&allFollowersID)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}

		var ids []uuid.UUID

		for i := 0; i < len(allFollowersID); i++ {
			ids = append(ids, allFollowersID[0].FollowerID)
		}

		res = tx.Where("id = ?", ids).Find(&followers)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return followers, err
}

func (r *Repo) GetAllFollowings(userID string) ([]entity.User, error) {
	var followings []entity.User
	err := r.replica.Db.Transaction(func(tx *gorm.DB) error {
		var allFollowingsID []entity.Follower
		res := tx.Where("follower_id = ?", userID).Find(&allFollowingsID)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}

		var ids []uuid.UUID

		for i := 0; i < len(allFollowingsID); i++ {
			ids = append(ids, allFollowingsID[0].FollowingID)
		}

		res = tx.Where("id = ?", ids).Find(&followings)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return followings, err
}
