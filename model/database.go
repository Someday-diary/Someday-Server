package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	AccessTokenRedis *redis.Client
	EmailVerifyRedis *redis.Client
)

func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("mysql_username"),
		os.Getenv("mysql_pwd"),
		os.Getenv("mysql_address"),
		os.Getenv("mysql_port"),
		os.Getenv("mysql_db_name"))

	sqlDB, err := sql.Open("mysql", dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(
		&User{},
		&Secret{},
		&Post{},
		&Tag{},
	)
	DB = db

	AccessTokenRedis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_address"),
		Password: os.Getenv("redis_pwd"), // no password set
		DB:       0,                      // use default DB
	})

	EmailVerifyRedis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_address"),
		Password: os.Getenv("redis_pwd"), // no password set
		DB:       1,                      // use default DB
	})

	ctx := context.Background()

	_, err = AccessTokenRedis.Ping(ctx).Result()
	if err != nil {
		log.Panic(err)
	}

	_, err = EmailVerifyRedis.Ping(ctx).Result()
	if err != nil {
		log.Panic(err)
	}

	log.Print("[DATABASE] 연결 완료")
}
