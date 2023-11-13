package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
)

func (r *GistRepo) GetOtherAllGists(sort string, direction string) (*[]entity.GistRequest, error) {
	var allGistsReq []entity.GistRequest

	var allGists []entity.Gist

	err := r.replica.Db.Order(fmt.Sprintf("%s %s", sort, direction)).Find(&allGists).Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(allGists); i++ {
		var lastCommit entity.Commit
		err = r.replica.Db.Where("gist_id = ?", allGists[i].ID).Order("created_at desc").Find(&lastCommit).Limit(1).Error
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

	return &allGistsReq, nil
}

func (r *GistRepo) GetGistByID(gistID uuid.UUID) (*entity.GistRequest, error) {
	var gistReq entity.GistRequest

	var gist entity.Gist

	res := r.replica.Db.Where("id = ?", gistID).Find(&gist)
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

	gistReq = entity.GistRequest{
		Gist:   gist,
		Commit: lastCommit,
		Files:  allFiles,
	}

	return &gistReq, nil
}

func (r *GistRepo) DeleteGistByID(id uuid.UUID) error {
	err := r.main.Db.Where("id = ?", id).Delete(&entity.Gist{}).Error
	return err
}
