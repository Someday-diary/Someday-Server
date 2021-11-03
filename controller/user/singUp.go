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
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		var user model.User
		var n int64
		result := model.DB.Where("email = ?", req.Email).First(&user).Count(&n)
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
			_ = c.Error(err)
			return
		}

		model.DB.Model(&user).Where("email = ?", req.Email).
			Updates(model.User{Pwd: sql.NullString{
				String: string(h),
				Valid:  true,
			}, Agree: req.Agree, Status: "normal"})

		key := lib.CreateToken(32)

		t, err := lib.SystemCipher.Encrypt(key)
		if err != nil {
			log.Panic(err)
		}

		//t := lib.EncryptAES([]byte(os.Getenv("secret_key")), key)

		model.DB.Create(&model.Secret{
			Email:     req.Email,
			SecretKey: t,
		})

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
