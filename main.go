package main

import (
	"github.com/Someday-diary/Someday-Server/controller"
	"github.com/Someday-diary/Someday-Server/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("", controller.SignUp)
		user.POST("/login", controller.Login)
		user.POST("/verify", controller.SendVerityEmail)
		user.POST("/verify/confirm", controller.Verity)
	}

	diary := r.Group("/diaries").Use(middleware.Auth)
	{
		diary.POST("", controller.CreatePost)
		diary.GET("", controller.GetPost)
		diary.GET("/:post_id", controller.GetPostByID)
		diary.PATCH("/:post_id", controller.EditPost)
		diary.DELETE("/:post_id", controller.DeletePost)
	}

	_ = r.Run()
}
