package model

import (
	"time"
)

type Post struct {
	ID        string    `gorm:"primaryKey" json:"id,omitempty"`
	Email     string    `json:"email" json:"email,omitempty"`
	Contents  string    `json:"contents" json:"contents,omitempty"`
	CreatedAt time.Time `json:"created_at" json:"created_at"`

	Tag []Tag `gorm:"foreignKey: PostID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"tag" json:"tag,omitempty"`
}

func (Post) TableName() string {
	return "post"
}
