package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

func (r *Repo) ForkGist(newFork entity.Fork) error {
	err := r.main.Db.Create(&newFork).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetForkedGistByUser(userID uuid.UUID) ([]entity.GistRequest, error) {
	var forks []entity.Fork

	err := r.replica.Db.Where("user_id = ?", userID).Find(&forks).Error
	if err != nil {
		return nil, err
	}

	var allGists []entity.GistRequest

	for i := 0; i < len(forks); i++ {
		var gist entity.Gist

		res := r.replica.Db.Where("id = ?", forks[i].GistID).Find(&gist)
		if res.Error != nil {
			return nil, res.Error
		}

		var lastCommit entity.Commit
		err := r.replica.Db.Where("gist_id = ?", gist.ID).Order("created_at desc").Find(&lastCommit).Limit(1).Error
		if err != nil {
			return nil, err
		}
		var allFiles []entity.File
		err = r.replica.Db.Where("commit_id = ?", lastCommit.ID).Find(&allFiles).Error
		if err != nil {
			return nil, err
		}

		gistReq := entity.GistRequest{
			Gist:   gist,
			Commit: lastCommit,
			Files:  allFiles,
		}
		allGists = append(allGists, gistReq)
	}
	return allGists, nil
}
