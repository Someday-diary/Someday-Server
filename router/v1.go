package router

import (
	"github.com/Someday-diary/Someday-Server/controller/v1/post"
	"github.com/Someday-diary/Someday-Server/controller/v1/server"
	"github.com/Someday-diary/Someday-Server/controller/v1/user"
	"github.com/Someday-diary/Someday-Server/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpV1() *gin.Engine {
	r := gin.Default()

	r.HEAD("/check", server.Check())
	r.HEAD("/error", server.Error())

	userAPI := r.Group("/user")
	{
		userAPI.POST("/sign_up", user.SignUp())
		userAPI.POST("/login", user.Login())
		userAPI.DELETE("/logout", middleware.Auth(), user.Logout())
		userAPI.POST("/verify", user.SendEmail())
		userAPI.POST("/verify/confirm", user.EmailConfirm())
	}

	postAPI := r.Group("/diary").Use(middleware.Auth())
	{
		postAPI.POST("", post.CreatePost())
		postAPI.GET("", post.GetPost())
		postAPI.GET("/month", post.GetPostByMonth())
		postAPI.GET("/date", post.GetPostByDate())
		postAPI.GET("/:post_id", post.GetPostByID())
		postAPI.PATCH("/:post_id", post.EditPost())
		postAPI.DELETE("/:post_id", post.DeletePost())
	}

	return r
}
