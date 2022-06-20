package dto

type Tag struct {
	TagName string `json:"tag_name"`
}

type Post struct {
	PostID   string `json:"post_id,omitempty"`
	Contents string `json:"contents,omitempty"`
	Date     string `json:"date,omitempty"`

	Tags *[]Tag `json:"tags,omitempty"`
}
