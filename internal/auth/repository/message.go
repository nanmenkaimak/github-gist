package repository

import (
	"fmt"
	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
)

func (r *Repo) InsertMessage(newMessage entity.Message) error {
	tx := r.main.Db.Create(&newMessage)
	if tx.Error != nil {
		return fmt.Errorf("failed inserting kafka message err: %v", tx.Error)
	}
	return nil
}

func (r *Repo) GetNotProcessedMessages() ([]entity.Message, error) {
	var allNotProcessedMessages []entity.Message
	err := r.replica.Db.Where("is_processed = ?", false).Find(&allNotProcessedMessages).Error
	if err != nil {
		return nil, err
	}
	return allNotProcessedMessages, nil
}

func (r *Repo) GetProcessedMessages() ([]entity.Message, error) {
	var allNotProcessedMessages []entity.Message
	err := r.replica.Db.Where("is_processed = ?", true).Find(&allNotProcessedMessages).Error
	if err != nil {
		return nil, err
	}
	return allNotProcessedMessages, nil
}

func (r *Repo) UpdateMessage(key string) error {
	tx := r.main.Db.Model(&entity.Message{}).Where("key = ?", key).Update("is_processed", true)
	if tx.Error != nil {
		return fmt.Errorf("failed updating kafka message err: %v", tx.Error)
	}
	return nil
}

func (r *Repo) DeleteMessage() error {
	err := r.main.Db.Where("is_processed = ?", true).Delete(&entity.Message{}).Error
	return err
}
