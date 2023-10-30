package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

func (r *Repo) CreateComment(newComment entity.Comment) error {
	err := r.main.Db.Create(&newComment).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetAllCommentsOfGist(gistID uuid.UUID) ([]entity.Comment, error) {
	var comments []entity.Comment

	err := r.replica.Db.Where("gist_id = ?", gistID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *Repo) DeleteComment(id uuid.UUID) error {
	err := r.main.Db.Where("id = ?", id).Delete(&entity.Comment{}).Error
	return err
}

func (r *Repo) UpdateComment(updatedComment entity.Comment) error {
	err := r.main.Db.Model(&updatedComment).Where("id = ?", updatedComment.ID).Updates(updatedComment).Error
	return err
}
