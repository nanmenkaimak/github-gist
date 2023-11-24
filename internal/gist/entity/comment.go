package entity

import (
	"time"

	"github.com/google/uuid"
)

// swagger:model Comment
type Comment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid; default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid; not null"`
	GistID    uuid.UUID `json:"gist_id" gorm:"not null"`
	Gist      Gist      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Text      string    `json:"text" gorm:"type:varchar(250); not null"`
	Username  string    `json:"username" gorm:"type:varchar(50); not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentBuilder struct {
	comment *Comment
}

func NewCommentBuilder() *CommentBuilder {
	return &CommentBuilder{
		comment: &Comment{},
	}
}

func (g *CommentBuilder) SetText(text string) *CommentBuilder {
	g.comment.Text = text
	return g
}

func (g *CommentBuilder) SetGistID(gistID uuid.UUID) *CommentBuilder {
	g.comment.GistID = gistID
	return g
}

func (g *CommentBuilder) SetUserID(userID uuid.UUID) *CommentBuilder {
	g.comment.UserID = userID
	return g
}

func (g *CommentBuilder) Build() *Comment {
	return g.comment
}
