package dto

type CreatePost struct {
	Tags []struct {
		TagName string `json:"tag" binding:"required"`
	} `json:"tags"`
	Contents string `json:"contents" binding:"required"`
	Date     string `json:"date" binding:"required"`
	ID       string `json:"id" binding:"required"`
}

type EditPost struct {
	Contents string `json:"contents" binding:"required"`
	Tags     []struct {
		TagName string `json:"tag" binding:"required"`
	} `json:"tags" binding:"required"`
}

type EmailConfirm struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type Login struct {
	Email string `json:"email" binding:"required"`
	Pwd   string `json:"pwd" binding:"required"`
}

type SendEmail struct {
	Email string `json:"email"`
}

type SignUp struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	Agree string `json:"agree"`
}
