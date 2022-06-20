package dao

type Tag struct {
	PostID  string `json:"post_id"`
	TagName string `json:"tag_name"`
}

func (Tag) TableName() string {
	return "tag"
}
