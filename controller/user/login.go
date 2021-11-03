package user

import (
	"context"
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Pwd   string `json:"pwd" binding:"required"`
}

func Login() gin.HandlerFunc {
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
		var count int64
		model.DB.Find(&user, "email = ?", req.Email).Count(&count)

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

		var k model.Secret
		model.DB.Select("secret_key").Where("email = ?", req.Email).Limit(1).Find(&k)
		secretKey, err := lib.SystemCipher.Decrypt(k.SecretKey)
		if err != nil {
			panic(err)
		}

		model.AccessTokenRedis.Set(context.Background(), token, req.Email, 0)

		c.JSON(http.StatusOK, gin.H{
			"code":       200,
			"token":      token,
			"secret_key": secretKey,
		})
	}
}
