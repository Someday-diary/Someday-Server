package user

import (
	"context"
	"net/http"

	"github.com/Someday-diary/Someday-Server/model"
	"github.com/gin-gonic/gin"
)

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := model.AccessTokenRedis.Del(context.Background(), c.GetHeader("access_token")).Result()
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
