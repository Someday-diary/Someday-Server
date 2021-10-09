package controller

import (
	"net/http"

	"github.com/Someday-diary/Someday-Server/lib"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	c.JSON(200, "call SI")
}

func Login(c *gin.Context) {
	c.JSON(200, "call LI")
}

func EmailVerification(c *gin.Context) {
	err := lib.SendEmail()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
