package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

func (r *Repo) StarGist(newStar entity.Star) error {
	err := r.main.Db.Create(&newStar).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetStarredGists(userID uuid.UUID) ([]entity.GistRequest, error) {
	var stars []entity.Star

	err := r.replica.Db.Where("user_id = ?", userID).Find(&stars).Error
	if err != nil {
		return nil, err
	}

	var allGists []entity.GistRequest

	for i := 0; i < len(stars); i++ {
		var gist entity.Gist

		res := r.replica.Db.Where("id = ?", stars[i].GistID).Find(&gist)
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

func (r *Repo) GetAllStargazers() {

}

func (r *Repo) DeleteStar(gistID uuid.UUID, userID uuid.UUID) error {
	err := r.main.Db.Where("gist_id = ? and user_id = ?", gistID, userID).Delete(&entity.Star{}).Error
	return err
}
