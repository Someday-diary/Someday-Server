package main

import (
	"log"
	"os"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/Someday-diary/Someday-Server/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("critical error! : %s", err.Error())
	}

	model.Connect()
	lib.SystemCipher = lib.CreateCipher(os.Getenv("secret_key"))

	r := router.SetUpV1()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("critical error! : %s", err.Error())
	}
}
