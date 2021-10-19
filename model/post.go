package model

import (
	"time"
)

type Post struct {
	ID        string `gorm:"primaryKey"`
	Email     string
	Contents  string
	CreatedAt time.Time

	Tag []Tag `gorm:"foreignKey: PostID"`
}

func (Post) TableName() string {
	return "post"
}
