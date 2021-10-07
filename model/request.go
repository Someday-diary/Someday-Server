package model

import "encoding/json"

type SignRequest struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type CreatePostRequest struct {
	Tags []struct {
		TagName string `json:"tag_name"`
	} `json:"tags"`
	Contents string `json:"contents"`
	Data     string `json:"data"`
}

type EditPostRequest struct {
	Tags []struct {
		TagName string `json:"tag_name"`
	} `json:"tags"`
	Contents string `json:"contents"`
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
