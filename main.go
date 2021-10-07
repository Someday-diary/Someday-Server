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
		user.POST("/signup", controller.SignUp)
		user.GET("/login", controller.Login)
	}

	diary := r.Group("/diary").Use(middleware.Auth)
	{
		diary.POST("", controller.CreatePost)
		diary.GET("", controller.GetAllPost)
		diary.GET("/:post_id", controller.GetPostByID)
		diary.GET("/search", controller.GetPostByTag)
		diary.PATCH("/:post_id", controller.EditPost)
		diary.DELETE("/:post_id", controller.DeletePost)
	}

	_ = r.Run()
}
