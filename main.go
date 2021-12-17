package main

import (
	"log"
	"os"

	"github.com/Someday-diary/Someday-Server/controller/post"
	"github.com/Someday-diary/Someday-Server/controller/user"
	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/middleware"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("critical error! : %s", err.Error())

	}

	model.Connect()
	lib.SystemCipher = lib.CreateCipher(os.Getenv("secret_key"))

	r := gin.Default()

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

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("critical error! : %s", err.Error())
	}
}
