package user

import (
	"context"
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
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
			c.JSON(http.StatusNotAcceptable, gin.H{
				"msg": err.Error(),
			})
			return
		}

		code, err := model.EmailVerifyRedis.Get(context.Background(), req.Email).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "유효한 인증 정보가 없습니다",
			})
			return
		}

		if code != req.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "입력한 코드가 맞지 않습니다",
			})
			return
		}

		var user model.User

		model.EmailVerifyRedis.Del(context.Background(), req.Email)
		model.DB.Model(&user).Select("status").Where("email = ?", req.Email).
			Updates(model.User{Status: "authenticated"})

		c.JSON(http.StatusOK, gin.H{
			"msg": "이메일 인증에 성공하였습니다!",
		})
	}
}
