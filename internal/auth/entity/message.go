package entity

import "github.com/google/uuid"

type Message struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Code        string    `json:"code" gorm:"type:char(4); not null"`
	Key         string    `json:"key" gorm:"type:varchar(100); not null"`
	IsProcessed *bool     `json:"is_processed" gorm:"default:false"`
}
