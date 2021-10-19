package main

import (
	"github.com/Someday-diary/Someday-Server/controller"
	"github.com/Someday-diary/Someday-Server/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.ErrorHandle())

	userAPI := r.Group("/user")
	{
		userAPI.POST("", user.SignUp())
		userAPI.POST("/login", user.Login())
		userAPI.POST("/verify", user.SendEmail())
		userAPI.POST("/verify/confirm", user.EmailConfirm())
	}

	postAPI := r.Group("/diary").Use(middleware.Auth())
	{
		postAPI.POST("", post.CreatePost())
		postAPI.GET("", post.GetPost())
		postAPI.GET("/:post_id", post.GetPostByID())
		postAPI.PATCH("/:post_id", post.EditPost())
		postAPI.DELETE("/:post_id", post.DeletePost())
	}

	_ = r.Run()
}
