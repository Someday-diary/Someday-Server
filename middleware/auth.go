package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/Someday-diary/Someday-Server/model/dao"
	"github.com/Someday-diary/Someday-Server/model/database"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.GetHeader("access_token")

		email, err := database.TokenDB.Get(context.Background(), t).Result()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
			})
			c.Abort()
			return
		}

		var secret dao.Secret
		database.DB.First(&secret, "email = ?", email)

		c.Request.Header.Add("email", email)

		ci, err := lib.NewNiceCrypto(os.Getenv("secret_key"), os.Getenv("cipher_iv_key"))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		k, err := ci.Decrypt(secret.SecretKey)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Request.Header.Set("secret_key", k)
		c.Next()
	}
}
