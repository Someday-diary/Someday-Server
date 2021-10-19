package model

import (
	"database/sql"
	"time"
)

type User struct {
	Email     string `gorm:"primaryKey"`
	Pwd       sql.NullString
	Agree     string `gorm:"type:enum('Y', 'N'); default:N;"`
	Status    string `gorm:"type:enum('normal','not authenticated', 'authenticated'); default:'not authenticated';"`
	CreatedAt time.Time

	Secret Secret `gorm:"foreignKey: Email"`
	Post   []Post `gorm:"foreignKey: Email"`
}

func (User) TableName() string {
	return "user"
}
