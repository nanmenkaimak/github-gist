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

func (r *Repo) GetForkedGistByUser(userID uuid.UUID, ownGists bool) ([]entity.GistRequest, error) {
	var allGistsReq []entity.GistRequest

	var allGists []entity.Gist

	gists := r.replica.Db.Where("user_id = ? and is_forked = true", userID)
	if !ownGists {
		gists.Where("visible = true")
	}
	gists.Find(&allGists)
	if gists.Error != nil {
		return nil, gists.Error
	}

	for i := 0; i < len(allGists); i++ {
		var lastCommit entity.Commit
		err := r.replica.Db.Where("gist_id = ?", allGists[i].ID).Order("created_at desc").Find(&lastCommit).Limit(1).Error
		if err != nil {
			return nil, err
		}
		var allFiles []entity.File
		err = r.replica.Db.Where("commit_id = ?", lastCommit.ID).Find(&allFiles).Error
		if err != nil {
			return nil, err
		}

		res := entity.GistRequest{
			Gist:   allGists[i],
			Commit: lastCommit,
			Files:  allFiles,
		}
		allGistsReq = append(allGistsReq, res)
	}

	return allGistsReq, nil
}

func (r *Repo) DeleteFork(id uuid.UUID) error {
	err := r.main.Db.Where("id = ?", id).Delete(&entity.Fork{}).Error
	return err
}
