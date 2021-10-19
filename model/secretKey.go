package model

type Secret struct {
	Email     string `gorm:"primaryKey"`
	SecretKey string
}

func (Secret) TableName() string {
	return "secret"
}
