package controller

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model/dao"
	"github.com/Someday-diary/Someday-Server/model/database"
	"github.com/Someday-diary/Someday-Server/model/dto"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmailConfirmRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func EmailConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(EmailConfirmRequest)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}
		code, err := database.EmailDB.Get(context.Background(), req.Email).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 102,
			})
			return
		}

		if code != req.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 103,
			})
			return
		}

		var user dao.User

		_, err = database.EmailDB.Del(context.Background(), req.Email).Result()
		if err != nil {
			panic(err)
		}
		err = database.DB.Model(&user).Select("status").Where("email = ?", req.Email).
			Updates(dao.User{Status: "authenticated"}).Error
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(dto.SignUp)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		var user dao.User
		var count int64
		database.DB.Find(&user, "email = ?", req.Email).Count(&count)

		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 106,
			})
			return
		}
		if user.Status == "not authenticated" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 107,
			})
			return
		}

		compare := bcrypt.CompareHashAndPassword([]byte(user.Pwd.String), []byte(req.Pwd))
		if compare != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 108,
			})
			return
		}

		token := lib.CreateToken(8)

		var k dao.Secret
		database.DB.Select("secret_key").Where("email = ?", req.Email).Limit(1).Find(&k)
		secretKey, err := lib.SystemCipher.Decrypt(k.SecretKey)
		if err != nil {
			panic(err)
		}

		database.TokenDB.Set(context.Background(), token, req.Email, 0)

		c.JSON(http.StatusOK, gin.H{
			"code":       200,
			"token":      token,
			"secret_key": secretKey,
		})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := database.TokenDB.Del(context.Background(), c.GetHeader("access_token")).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 112,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}

func SendEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(dto.SendEmail)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		var u dao.User
		var n int64
		err = database.DB.Find(&u, "email = ?", req.Email).Count(&n).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}

		if n >= 1 {
			if u.Status == "normal" {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 101,
				})
				return
			}
		}

		u = dao.User{
			Email: req.Email,
			Pwd: sql.NullString{
				String: "",
				Valid:  false,
			},
			Agree:     "N",
			Status:    "not authenticated",
			CreatedAt: time.Now(),
		}

		err = database.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&u).Error
		if err != nil {
			panic(err)
		}

		code, err := lib.CreateCode()
		if err != nil {
			panic(err)
		}
		database.TokenDB.Set(context.Background(), req.Email, code, time.Minute*30)

		templateData := map[string]string{
			"code": code,
		}

		r := lib.NewRequest(req.Email, os.Getenv("smtp_id"), "[오늘하루] 회원가입을 위한 인증번호를 알려드립니다", "")
		err = r.ParseTemplate("templates/email.html", templateData)
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

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(dto.SignUp)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		var user dao.User
		var n int64
		result := database.DB.Where("email = ?", req.Email).First(&user).Count(&n)
		if result.Error != nil || user.Status == "not authenticated" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 104,
			})
			return
		}

		if user.Status == "normal" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 105,
			})
			return
		}

		h, err := bcrypt.GenerateFromPassword([]byte(req.Pwd), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		database.DB.Model(&user).Where("email = ?", req.Email).
			Updates(dao.User{Pwd: sql.NullString{
				String: string(h),
				Valid:  true,
			}, Agree: req.Agree, Status: "normal"})

		key := lib.CreateToken(32)

		t, err := lib.SystemCipher.Encrypt(key)
		if err != nil {
			log.Panic(err)
		}

		database.DB.Create(&dao.Secret{
			Email:     req.Email,
			SecretKey: t,
		})

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
