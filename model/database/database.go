package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Someday-diary/Someday-Server/model/dao"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm/logger"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("mysql_username"),
		os.Getenv("mysql_pwd"),
		os.Getenv("mysql_address"),
		os.Getenv("mysql_port"),
		os.Getenv("mysql_db_name"))
	sqlDB, err := sql.Open("mysql", dsn)
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, DefaultStringSize: 191}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Panic(err)
	}

	err = db.AutoMigrate(
		&dao.User{},
		&dao.Secret{},
		&dao.Post{},
		&dao.Tag{},
	)
	if err != nil {
		log.Panic(err)
	}

	log.Print("[DATABASE] 연결 완료")
}

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		os.Getenv("mysql_username"),
		os.Getenv("mysql_pwd"),
		os.Getenv("mysql_address"),
		os.Getenv("mysql_port"),
		os.Getenv("mysql_db_name"))

	sqlDB, err := sql.Open("mysql", dsn)
	db, err := gorm.Open(mysql.New(
		mysql.Config{
			Conn:              sqlDB,
			DefaultStringSize: 191}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Panic(err)
	}
	return db
}

func ConnectTokenRedis() *redis.Client {
	ctx := context.Background()
	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_address"),
		Password: os.Getenv("redis_pwd"), // no password set
		DB:       0,                      // use default DB
	})

	_, err := r.Ping(ctx).Result()
	if err != nil {
		log.Panic(err)
	}
	return r
}

func ConnectEmailRedis() *redis.Client {
	ctx := context.Background()
	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_address"),
		Password: os.Getenv("redis_pwd"), // no password set
		DB:       1,                      // use default DB
	})

	_, err := r.Ping(ctx).Result()
	if err != nil {
		log.Panic(err)
	}
	return r
}
