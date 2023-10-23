package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"gorm.io/gorm"
	"time"
)

func (r *Repo) CreateGist(newGist entity.GistRequest) (uuid.UUID, error) {
	err := r.main.Db.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(&newGist.Gist)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		newGist.Commit.GistID = newGist.Gist.ID
		res = tx.Create(&newGist.Commit)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}

		for i := 0; i < len(newGist.Files); i++ {
			newGist.Files[i].CommitID = newGist.Commit.ID
		}

		res = tx.Create(&newGist.Files)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return newGist.Gist.ID, nil
}

func (r *Repo) GetOtherAllGists() ([]entity.GistRequest, error) {
	var allGistsReq []entity.GistRequest

	var allGists []entity.Gist

	err := r.replica.Db.Where("visible = true").Find(&allGists).Error
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

	return allGistsReq, nil
}

func (r *Repo) GetGistByID(gistID uuid.UUID, ownGist bool) (entity.GistRequest, error) {
	var gistReq entity.GistRequest

	var gist entity.Gist

	res := r.replica.Db.Where("id = ?", gistID)
	if !ownGist {
		res.Where("visible = true")
	}
	res.Find(&gist)
	if res.Error != nil {
		return gistReq, res.Error
	}

	var lastCommit entity.Commit
	err := r.replica.Db.Where("gist_id = ?", gist.ID).Order("created_at desc").Find(&lastCommit).Limit(1).Error
	if err != nil {
		return gistReq, err
	}
	var allFiles []entity.File
	err = r.replica.Db.Where("commit_id = ?", lastCommit.ID).Find(&allFiles).Error
	if err != nil {
		return gistReq, err
	}

	gistReq = entity.GistRequest{
		Gist:   gist,
		Commit: lastCommit,
		Files:  allFiles,
	}

	return gistReq, nil
}

func (r *Repo) GetAllGistsOfUser(userID uuid.UUID, ownGists bool) ([]entity.GistRequest, error) {
	var allGistsReq []entity.GistRequest

	var allGists []entity.Gist

	gists := r.replica.Db.Where("user_id = ?", userID)
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

func (r *Repo) UpdateGistByID(updatedGist entity.GistRequest) error {
	err := r.main.Db.Transaction(func(tx *gorm.DB) error {
		updatedGist.Gist.UpdatedAt = time.Now()
		res := tx.Model(&updatedGist.Gist).Where("id = ?", updatedGist.Gist.ID).Updates(updatedGist.Gist)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		updatedGist.Commit.GistID = updatedGist.Gist.ID
		res = tx.Create(&updatedGist.Commit)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}

		for i := 0; i < len(updatedGist.Files); i++ {
			updatedGist.Files[i].CommitID = updatedGist.Commit.ID
		}

		res = tx.Create(&updatedGist.Files)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
		return nil
	})
	return err
}

func (r *Repo) DeleteGistByID(id uuid.UUID) error {
	err := r.main.Db.Where("id = ?", id).Delete(&entity.Gist{}).Error
	return err
}
