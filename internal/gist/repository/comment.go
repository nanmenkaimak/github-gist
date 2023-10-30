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

func (r *Repo) UpdateComment() {
}
