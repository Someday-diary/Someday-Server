package router

import (
	con "github.com/Someday-diary/Someday-Server/controller"
	"github.com/Someday-diary/Someday-Server/middleware"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.HEAD("/check", con.HealthCheck())
	r.HEAD("/error", con.ReturnError())

	userAPI := r.Group("/user")
	{
		userAPI.POST("/sign_up", con.SignUp())
		userAPI.POST("/login", con.Login())
		userAPI.DELETE("/logout", middleware.Auth(), con.Logout())
		userAPI.POST("/verify", con.SendEmail())
		userAPI.POST("/verify/confirm", con.EmailConfirm())
	}

	postAPI := r.Group("/diary").Use(middleware.Auth())
	{
		postAPI.POST("", con.CreatePost())
		postAPI.GET("", con.GetPost())
		postAPI.GET("/month", con.GetPostByMonth())
		postAPI.GET("/date", con.GetPostByDate())
		postAPI.GET("/:post_id", con.GetPostByID())
		postAPI.PATCH("/:post_id", con.EditPost())
		postAPI.DELETE("/:post_id", con.DeletePost())
	}

	return r
}
