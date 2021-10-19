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
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

func Login() gin.HandlerFunc {
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
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "일치하는 유저가 없습니다.",
			})
			return
		}
		if user.Status == "not authenticated" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "이메일 인증을 하지 않았습니다.",
			})
			return
		}

		compare := bcrypt.CompareHashAndPassword([]byte(user.Pwd.String), []byte(req.Pwd))
		if compare != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "비밀번호가 일치하지 않습니다.",
			})
			return
		}

		token := lib.CreateToken(8)

		var k model.Secret
		model.DB.Select("secret_key").Where("email = ?", req.Email).First(&k)
		secretKey, err := lib.Cipher.Decrypt(k.SecretKey)
		if err != nil {
			_ = c.Error(err)
			return
		}

		model.AccessTokenRedis.Set(context.Background(), token, req.Email, 0)

		c.JSON(http.StatusOK, gin.H{
			"token":      token,
			"secret_key": secretKey,
		})
	}
}
