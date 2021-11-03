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
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
			})
			return
		}

		code, err := model.EmailVerifyRedis.Get(context.Background(), req.Email).Result()
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

		var user model.User

		model.EmailVerifyRedis.Del(context.Background(), req.Email)
		model.DB.Model(&user).Select("status").Where("email = ?", req.Email).
			Updates(model.User{Status: "authenticated"})

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
		})
	}
}
