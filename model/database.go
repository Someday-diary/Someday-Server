package model

import (
	"context"
	_ "database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectMYSQL() (*sqlx.DB, error) {
	return sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/someday", os.Getenv("mysql_username"),
		os.Getenv("mysql_pwd"), os.Getenv("mysql_address")))
}

func ConnectRedis(dbCode int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_address"),
		Password: os.Getenv("redis_pwd"), // no password set
		DB:       dbCode,                 // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	return client, err
}
