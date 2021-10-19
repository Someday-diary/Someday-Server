package user

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
)

type SendEmailRequest struct {
	Email string `json:"email"`
}

func SendEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(SendEmailRequest)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"msg": err.Error(),
			})
			return
		}

		err = model.DB.Select("email").
			Create(&model.User{Email: req.Email, CreatedAt: time.Now().Add(time.Hour * 9)}).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "이미 가입되어있는 유저입니다.",
			})
			return
		}

		code, err := lib.CreateCode()
		if err != nil {
			_ = c.Error(err)
			return
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
			_ = c.Error(err)
			return
		}
	}
}
