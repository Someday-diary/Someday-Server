package dao

import (
	"database/sql"
	"time"
)

type User struct {
	Email     string         `gorm:"primaryKey" json:"email"`
	Pwd       sql.NullString `json:"pwd"`
	Agree     string         `gorm:"type:enum('Y', 'N'); default:N;"`
	Status    string         `gorm:"type:enum('normal','not authenticated', 'authenticated'); default:'not authenticated';" json:"status"`
	CreatedAt time.Time      `json:"created_at"`

	Secret Secret `gorm:"foreignKey: Email; constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"secret"`
	Post   []Post `gorm:"foreignKey: Email; constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"post"`
}

func (User) TableName() string {
	return "user"
}
