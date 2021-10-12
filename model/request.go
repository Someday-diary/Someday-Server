package model

import "encoding/json"

type SignRequest struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	Agree string `json:"agree"`
}

type CreatePostRequest struct {
	Diaries []struct {
		Tags []struct {
			TagName string `json:"tag_name" db:"tag_name"`
		} `json:"tags"`
		Contents string `json:"contents" db:"contents"`
		Date     string `json:"data" db:"created_at"`
	}
}

type EditPostRequest struct {
	Tags []struct {
		TagName string `json:"tag_name"`
	} `json:"tags"`
	Contents string `json:"contents"`
}

type EmailVerityRequest struct {
	Email string `json:"email"`
}

func UnmarshalSignRequest(data []byte) (SignRequest, error) {
	var r SignRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalCreatePostRequest(data []byte) (CreatePostRequest, error) {
	var r CreatePostRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalEmailVerityRequest(data []byte) (EmailVerityRequest, error) {
	var r EmailVerityRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func UnmarshalEditPostRequest(data []byte) (EditPostRequest, error) {
	var r EditPostRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SignRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *CreatePostRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *EditPostRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
