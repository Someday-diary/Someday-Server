package main

import (
	"log"
	"os"

	"github.com/Someday-diary/Someday-Server/lib"
	_ "github.com/Someday-diary/Someday-Server/model/database"
	"github.com/Someday-diary/Someday-Server/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("critical error! : %s", err.Error())
	}

	lib.SystemCipher = lib.CreateCipher(os.Getenv("secret_key"))

	r := router.Setup()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("critical error! : %s", err.Error())
	}
}
