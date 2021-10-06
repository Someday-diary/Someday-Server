package main

import (
	"github.com/Someday-diary/Someday-Server/Controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/sign-in", Controller.SignIn)
		user.GET("/login", Controller.Login)
	}

	diary := r.Group("/diary")
	{
		diary.POST("", Controller.CreatePost)
		diary.GET("", Controller.GetAllPost)
		diary.GET("/:post_id", Controller.GetPostByID)
		diary.GET("/search", Controller.GetPostByTag)
		diary.PATCH("/:post_id", Controller.EditPost)
		diary.DELETE("/:post_id", Controller.DeletePost)
	}

	_ = r.Run()
}
