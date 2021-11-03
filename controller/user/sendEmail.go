package user

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SendEmailRequest struct {
	Email string `json:"email"`
}

func SendEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(SendEmailRequest)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		var u model.User
		var n int64
		err = model.DB.Find(&u, "email = ?", req.Email).Count(&n).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}

		if n >= 1 {
			if u.Status == "normal" || u.Status == "authenticated" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 101,
				})
				return
			}
		}

		u = model.User{
			Email: req.Email,
			Pwd: sql.NullString{
				String: "",
				Valid:  false,
			},
			Agree:     "N",
			Status:    "not authenticated",
			CreatedAt: time.Now(),
		}

		err = model.DB.Create(&u).Error
		if err != nil {
			panic(err)
		}

		code, err := lib.CreateCode()
		if err != nil {
			panic(err)
		}
		model.EmailVerifyRedis.Set(context.Background(), req.Email, code, time.Minute*30)

		templateData := map[string]string{
			"code": code,
		}

		r := lib.NewRequest(req.Email, os.Getenv("smtp_id"), "[오늘하루] 회원가입을 위한 인증번호를 알려드립니다", "")
		err = r.ParseTemplate("templates/index.html", templateData)
		if err != nil {
			log.Panic(err)
		}

		err = r.SendEmail()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
