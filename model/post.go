package model

import (
	"time"
)

type Post struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Email     string    `json:"email"`
	Contents  string    `json:"contents"`
	CreatedAt time.Time `json:"created_at"`

	Tag []Tag `gorm:"foreignKey: PostID" json:"tag"`
}

func (Post) TableName() string {
	return "post"
}
