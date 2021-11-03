package model

type Secret struct {
	Email     string `gorm:"primaryKey" json:"email"`
	SecretKey string `json:"secret_key"`
}

func (Secret) TableName() string {
	return "secret"
}
