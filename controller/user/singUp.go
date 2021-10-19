package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	Agree string `json:"agree"`
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := new(SignUpRequest)
		err := c.Bind(req)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"msg": err.Error(),
			})
			return
		}

		var user model.User
		result := model.DB.Where("email = ?", req.Email).First(&user)
		if result.Error != nil || user.Status == "not authenticated" {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "이메일 인증 안함 응애",
			})
			return
		}

		if user.Status == "normal" {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "이미 있는 계정입니다",
			})
			return
		}

		h, err := bcrypt.GenerateFromPassword([]byte(req.Pwd), bcrypt.DefaultCost)
		if err != nil {
			_ = c.Error(err)
			return
		}

		model.DB.Model(&user).Where("email = ?", req.Email).
			Updates(model.User{Pwd: sql.NullString{
				String: string(h),
				Valid:  true,
			}, Agree: req.Agree, Status: "normal"})

		key := lib.CreateToken(16)

		t, err := lib.Cipher.Encrypt(key)
		if err != nil {
			log.Panic(err)
		}

		//t := lib.EncryptAES([]byte(os.Getenv("secret_key")), key)

		model.DB.Create(&model.Secret{
			Email:     req.Email,
			SecretKey: t,
		})

		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	}
}
