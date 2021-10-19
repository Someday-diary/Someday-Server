package model

type Tag struct {
	PostID  string
	TagName string
}

func (Tag) TableName() string {
	return "tag"
}
